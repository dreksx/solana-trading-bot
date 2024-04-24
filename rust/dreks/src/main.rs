use futures_util::StreamExt;
use solana_client::nonblocking::pubsub_client::PubsubClient;

#[derive(serde::Deserialize)]
struct Env {
    ws_url: url::Url,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let env = envy::from_env::<Env>()?;

    let ps_client = PubsubClient::new(&env.ws_url.to_string()).await?;

    let program_id = Pubkey::from_str("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")?;
    let filters = vec![
            RpcAccountInfoFilter::DataSize(LIQUIDITY_STATE_LAYOUT_V4.span as usize),
            RpcAccountInfoFilter::Memcmp {
                offset: LIQUIDITY_STATE_LAYOUT_V4.offset_of("quoteMint")?,
                bytes: bs58::decode(quote_token_mint).into_vec()?,
            },
            RpcAccountInfoFilter::Memcmp {
                offset: LIQUIDITY_STATE_LAYOUT_V4.offset_of("marketProgramId")?,
                bytes: bs58::decode(mainnet_program_id_openbook_market).into_vec()?,
            },
            RpcAccountInfoFilter::Memcmp {
                offset: LIQUIDITY_STATE_LAYOUT_V4.offset_of("status")?,
                bytes: bs58::encode(&[6, 0, 0, 0, 0, 0, 0, 0]).into_vec()?,
            },
        ];

        let config = RpcProgramAccountsConfig {
            filters: Some(filters),
            account_config: RpcAccountInfoConfig {
                encoding: Some(String::from("base64")),
                commitment: Some(CommitmentConfig::confirmed()),
                ..Default::default()
            },
            ..Default::default()
        };

        let subscription = client.on_program_account_change(
            &Pubkey::from_str(mainnet_program_id_amm_v4)?,
            move |updated_account_info| {
                // Process updated account info here
                println!("Updated account: {:?}", updated_account_info);
            },
            &config,
        ).await?;

    let (mut accs, unsubscriber) = ps_client.slot_subscribe().await?;

    let mut count = 0;
    while let Some(response) = accs.next().await {
        println!("{:?}", response);
        count += 1;
        if count >= 5 {
            break;
        }
    }

    unsubscriber().await;

    Ok(())
}
