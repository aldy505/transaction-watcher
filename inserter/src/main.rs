use std::env;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::thread;
use std::time;
use rand::{thread_rng, Rng};
use postgres::{Client, Error, NoTls, Transaction};

fn main() {
    let database_url = env::var("DATABASE_URL").unwrap_or(String::from("postgresql://watcher:password@localhost:5432/watcher?sslmode=disable"));
    let mut client = Client::connect(database_url.as_str(), NoTls)?;
    let mut rng = thread_rng();

    let running = Arc::new(AtomicBool::new(true));
    let r = running.clone();

    ctrlc::set_handler(move || {
        r.store(false, Ordering::SeqCst);
        client.close()?;
    }).expect("Error setting Ctrl-C handler");

    while running.load(Ordering::SeqCst)  {
        let transaction_amount = rng.gen_range(-1e5..=1e9);
        let transaction_type = match rng.gen_range(0..=3) {
            0 => String::from("TOP_UP"),
            1 => String::from("TRANSFER"),
            2 => String::from("WITHDRAW"),
            3 => String::from("FEE"),
            _ => String::from("")
        };
        let customer_number = rng.gen_range(100..200);

        match client.transaction() {
            Ok(mut transaction) => {
                match transaction.execute(
                    "INSERT INTO transactions (transaction_type, customer_number, transaction_amount, timestamp) VALUES ($1, $2, $3, NOW())",
                    &[transaction_type, customer_number, transaction_amount]
                ) {
                    Ok(_) => {}
                    Err(e1) => {
                        eprintln!("{:?}", e1);
                        if let Err(e2) = transaction.rollback() {
                            eprintln!("{:?}", e2)
                        }
                    }
                }

                match transaction.commit() {
                    Ok(_) => {}
                    Err(e1) => {
                        eprintln!("{:?}", e1);
                        if let Err(e2) = transaction.rollback() {
                            eprintln!("{:?}", e2)
                        }
                    }
                }
            }
            Err(error) => {
                eprintf!("{:?}", error)
            }
        }

        let time_to_sleep = time::Duration::from_millis(rng.gen_range(10..1000));
        thread::sleep(time_to_sleep);
    }
}
