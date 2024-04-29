import { LIQUIDITY_STATE_LAYOUT_V4, MAINNET_PROGRAM_ID, MARKET_STATE_LAYOUT_V3, Token } from '@raydium-io/raydium-sdk';
import bs58 from 'bs58';
import { Connection, PublicKey } from '@solana/web3.js';
import { TOKEN_PROGRAM_ID } from '@solana/spl-token';
import { EventEmitter } from 'events';
import WebSocket from 'ws';

export class Listeners extends EventEmitter {
  private subscriptions: number[] = [];

  constructor(private readonly connection: Connection, private readonly connectionPremium: WebSocket) {
    super();
  }

  public async start(config: {
    walletPublicKey: PublicKey;
    quoteToken: Token;
    autoSell: boolean;
    cacheNewMarkets: boolean;
  }) {    
    //const transactionListener = await this.subscribeToTransactions(config);
    // if (config.cacheNewMarkets) {
    //   const openBookSubscription = await this.subscribeToOpenBookMarkets(config);
    //   this.subscriptions.push(openBookSubscription);
    // }

    const raydiumSubscription = await this.subscribeToRaydiumPools(config);
    this.subscriptions.push(raydiumSubscription);

    // if (config.autoSell) {
    //   const walletSubscription = await this.subscribeToWalletChanges(config);
    //   this.subscriptions.push(walletSubscription);
    // }
  }


  private async subscribeToTransactions(config: { quoteToken: Token }) {
    const connectionPremium = this.connectionPremium
    let sendRequest = function() {
      const request = {
        jsonrpc: "2.0",
        id: 420,
        method: "transactionSubscribe",
        params: [
            {
                accountRequired: ["675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8", "25hAyBQfoDhfWx9ay6rarbgvWGwDdNqcHsXS3jQ3mTDJ", "5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1"],
                accountExclude: ["JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4"]
            },
            {
                commitment: "processed",
                encoding: "jsonParsed",
                transactionDetails: "full",
                showRewards: true,
                maxSupportedTransactionVersion: 0
            }
        ]
    };
    connectionPremium.send(JSON.stringify(request));
  }

    this.connectionPremium.on("open", function open() {
      console.log('WebSocket is open');
      sendRequest();  // Send a request once the WebSocket is open
    });

    const self = this
    this.connectionPremium.on('message', function incoming(data: any) {
      const messageStr = data.toString('utf8');
      try {
          const messageObj = JSON.parse(messageStr);
          const transaction = messageObj?.params?.result?.transaction;
          if (transaction && transaction.meta.err === null) {
            self.emit('transaction', messageObj.params.result);
          }
      } catch (e) {
          console.error('Failed to parse JSON:', e);
      }
    });
  }

  private async subscribeToBlocks(config: { quoteToken: Token }) {
    return this.connection.onProgramAccountChange(
        MAINNET_PROGRAM_ID.OPENBOOK_MARKET,
        async (updatedAccountInfo) => {
          this.emit('market', updatedAccountInfo);
        },
        this.connection.commitment,
        [
          { dataSize: MARKET_STATE_LAYOUT_V3.span },
          {
            memcmp: {
              offset: MARKET_STATE_LAYOUT_V3.offsetOf('quoteMint'),
              bytes: config.quoteToken.mint.toBase58(),
            },
          },
        ],
    );
  }

  private async subscribeToOpenBookMarkets(config: { quoteToken: Token }) {
    return this.connection.onProgramAccountChange(
      MAINNET_PROGRAM_ID.OPENBOOK_MARKET,
      async (updatedAccountInfo) => {
        this.emit('market', updatedAccountInfo);
      },
      this.connection.commitment,
      [
        { dataSize: MARKET_STATE_LAYOUT_V3.span },
        {
          memcmp: {
            offset: MARKET_STATE_LAYOUT_V3.offsetOf('quoteMint'),
            bytes: config.quoteToken.mint.toBase58(),
          },
        },
      ],
    );
  }

  private async subscribeToRaydiumPools(config: { quoteToken: Token }) {
    console.log(LIQUIDITY_STATE_LAYOUT_V4.offsetOf('poolOpenTime'), LIQUIDITY_STATE_LAYOUT_V4.offsetOf('baseMint'))
    return this.connection.onProgramAccountChange(
      MAINNET_PROGRAM_ID.AmmV4,
      async (updatedAccountInfo) => {
        this.emit('pool', updatedAccountInfo);
      },
      this.connection.commitment,
      [
        { dataSize: LIQUIDITY_STATE_LAYOUT_V4.span },
        {
          memcmp: {
            offset: LIQUIDITY_STATE_LAYOUT_V4.offsetOf('quoteMint'),
            bytes: config.quoteToken.mint.toBase58(),
          },
        },
        {
          memcmp: {
            offset: LIQUIDITY_STATE_LAYOUT_V4.offsetOf('marketProgramId'),
            bytes: MAINNET_PROGRAM_ID.OPENBOOK_MARKET.toBase58(),
          },
        },
        {
          memcmp: {
            offset: LIQUIDITY_STATE_LAYOUT_V4.offsetOf('status'),
            bytes: bs58.encode([6, 0, 0, 0, 0, 0, 0, 0]),
          },
        },
      ],
    );
  }

  private async subscribeToWalletChanges(config: { walletPublicKey: PublicKey }) {
    return this.connection.onProgramAccountChange(
      TOKEN_PROGRAM_ID,
      async (updatedAccountInfo) => {
        this.emit('wallet', updatedAccountInfo);
      },
      this.connection.commitment,
      [
        {
          dataSize: 165,
        },
        {
          memcmp: {
            offset: 32,
            bytes: config.walletPublicKey.toBase58(),
          },
        },
      ],
    );
  }

  public async stop() {
    for (let i = this.subscriptions.length; i >= 0; --i) {
      const subscription = this.subscriptions[i];
      await this.connection.removeAccountChangeListener(subscription);
      this.subscriptions.splice(i, 1);
    }
  }
}
