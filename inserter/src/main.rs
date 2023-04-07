use std::env;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::thread;
use std::time;
use rand::{thread_rng, Rng};
use postgres::{Client, NoTls};

fn main() {
    let database_url = env::var("DATABASE_URL")
        .unwrap_or(String::from("postgresql://watcher:password@localhost:5432/watcher?sslmode=disable"));
    let mut client = Client::connect(database_url.as_str(), NoTls).unwrap();
    let mut rng = thread_rng();

    let running = Arc::new(AtomicBool::new(true));
    let r = running.clone();

    ctrlc::set_handler(move || {
        r.store(false, Ordering::SeqCst);
    }).expect("Error setting Ctrl-C handler");

    while running.load(Ordering::SeqCst) {
        let transaction_amount: i64 = rng.gen_range(0..1e9 as i64);
        let transaction_type = match rng.gen_range(0..=3) {
            0 => String::from("TOP_UP"),
            1 => String::from("TRANSFER"),
            2 => String::from("WITHDRAW"),
            3 => String::from("FEE"),
            _ => String::from("")
        };
        let customer_number = rng.gen_range(100..200);

        let exec_error = client.execute(
            "INSERT INTO transactions (transaction_type, customer_number, transaction_amount, timestamp) VALUES ($1, $2, $3, NOW())",
            &[&transaction_type, &customer_number, &transaction_amount]
        ).err();

        println!("{}", format!(
                "INSERT INTO transactions (transaction_type, customer_number, transaction_amount, timestamp) VALUES ({0}, {1}, {2}, NOW())",
                transaction_type,
                customer_number,
                transaction_amount
            ).to_string());

        match exec_error {
            None => {}
            Some(e) => {
                println!("{:?}", e);
            }
        }

        let time_to_sleep = time::Duration::from_millis(rng.gen_range(10..1000));
        thread::sleep(time_to_sleep);
    }

    client.close().unwrap();
}
