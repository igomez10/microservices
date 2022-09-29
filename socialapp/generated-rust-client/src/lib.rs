#[macro_use]
extern crate serde_derive;

use futures::executor::block_on;

extern crate reqwest;
extern crate serde;
extern crate serde_json;
extern crate tokio;
extern crate url;

pub mod apis;
pub mod models;
