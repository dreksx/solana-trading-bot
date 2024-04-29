use std::str::FromStr;
use futures_util::StreamExt;
use solana_sdk::pubkey::Pubkey;
use solana_client::nonblocking::pubsub_client::PubsubClient;
use solana_client::rpc_config::RpcProgramAccountsConfig;
use solana_account_decoder::{UiAccountEncoding, UiDataSliceConfig};
use solana_sdk::commitment_config::{CommitmentConfig, CommitmentLevel};
use solana_client::rpc_filter::{RpcFilterType, Memcmp, MemcmpEncodedBytes};
use byteorder::{ByteOrder, LittleEndian};
use std::io::Cursor;
use std::time::{SystemTime, UNIX_EPOCH};
use std::collections::HashMap;

#[derive(Debug)]
struct LiquidityStateV4 {
    status: u64,
    nonce: u64,
    max_order: u64,
    depth: u64,
    base_decimal: u64,
    quote_decimal: u64,
    state: u64,
    reset_flag: u64,
    min_size: u64,
    vol_max_cut_ratio: u64,
    amount_wave_ratio: u64,
    base_lot_size: u64,
    quote_lot_size: u64,
    min_price_multiplier: u64,
    max_price_multiplier: u64,
    system_decimal_value: u64,
    min_separate_numerator: u64,
    min_separate_denominator: u64,
    trade_fee_numerator: u64,
    trade_fee_denominator: u64,
    pnl_numerator: u64,
    pnl_denominator: u64,
    swap_fee_numerator: u64,
    swap_fee_denominator: u64,
    base_need_take_pnl: u64,
    quote_need_take_pnl: u64,
    quote_total_pnl: u64,
    base_total_pnl: u64,
    pool_open_time: u64,
    punish_pc_amount: u64,
    punish_coin_amount: u64,
    orderbook_to_init_time: u64,
    swap_base_in_amount: [u8; 16],
    swap_quote_out_amount: [u8; 16],
    swap_base2_quote_fee: u64,
    swap_quote_in_amount: [u8; 16],
    swap_base_out_amount: [u8; 16],
    swap_quote2_base_fee: u64,
    base_vault: Pubkey,
    quote_vault: Pubkey,
    base_mint: Pubkey,
    quote_mint: Pubkey,
}


#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let program = Pubkey::from_str("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8").unwrap();

    let ps_client = PubsubClient::new("wss://mainnet.helius-rpc.com/?api-key=b06c8c49-c031-4a59-8e6b-02ae50f63113").await?;



    let quote_mint = Pubkey::from_str("So11111111111111111111111111111111111111112").unwrap();
    let open_book = Pubkey::from_str("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX").unwrap();

    ;

    let vec = vec![6, 0, 0, 0, 0, 0, 0, 0];

    let filters = vec![
        RpcFilterType::Memcmp(Memcmp::new(432, MemcmpEncodedBytes::Base58(quote_mint.to_string()))),
        RpcFilterType::Memcmp(Memcmp::new(560, MemcmpEncodedBytes::Base58(open_book.to_string()))),
        RpcFilterType::Memcmp(Memcmp::new(0, MemcmpEncodedBytes::Bytes(vec))),
    ];
    let commitment = CommitmentConfig::processed();
    let config = RpcProgramAccountsConfig {
        filters: Some(filters),
        account_config: solana_client::rpc_config::RpcAccountInfoConfig {
            encoding: Some(UiAccountEncoding::Base64),
            data_slice: None,
            commitment: Some(commitment),
            min_context_slot: None,
        },
        with_context: None,
    };

    let (mut accounts, unsubscriber) = ps_client.
        program_subscribe(&program, Option::from(config)).await?;

    let mut count = 0;
    let mut now= SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() as u32;
    let mut cache = HashMap::new();

    while let Some(response) = accounts.next().await {
        let decoded = response.value.account.data.decode().unwrap();

        let mut openTime = LittleEndian::read_u32(&decoded[224..224+8]);
        if openTime < now {
            continue;
        }
        if decoded.len() >= 432 {
            let mintBytes: [u8;32] = decoded[400..432].try_into().expect("Slice with incorrect length");
            let mintValue = Pubkey::from(mintBytes);
            let key = mintValue.to_string();
            println!("{:?} {:?}, {:?}", SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_millis(), openTime, mintValue);
            if cache.contains_key(&key) {
                continue
            }

            cache.insert(key, 1);
        }
    }

    unsubscriber().await;

    Ok(())
}