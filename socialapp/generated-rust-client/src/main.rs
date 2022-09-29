// import datetime

// use apis::configuration::Configuration;
// // Struct
// struct Point {
//     x: i32,
//     y: i32,
// }

// struct User {
//     first_name: String,
//     last_name: String,
//     email: Option<String>,
//     age: i32,
// }

// fn main() {
//     println!("this is the message: {}", add_two_numbers(1, 2));
//     //list of names
//     let names = vec!["John", "Jane", "Jack"];
//     // print each name
//     for name in names {
//         println!("Hello {}", name);
//     }
//     // create a point
//     let point = Point { x: 1, y: 2 };
//     println!("point: {}", point);

//     // list of users
//     let users = vec![
//         User {
//             first_name: "John".to_string(),
//             last_name: "Doe".to_string(),
//             email: Some("myemail".to_string()),
//             age: 20,
//         },
//         User {
//             first_name: "Jane".to_string(),
//             last_name: "Doe".to_string(),
//             email: None,
//             age: 20,
//         },
//     ];

//     // print each user
//     for user in users {
//         println!("Full name: {}", user.full_name());
//         println!("\tage: {}", user.age());
//         if user.email != None {
//             println!("\temail: {}", user.email.unwrap());
//         }
//     }
// }

// // add_two_numbers adds two numbers together
// fn add_two_numbers(a: i32, b: i32) -> i32 {
//     return a + b;
// }

// // add function display to point
// impl std::fmt::Display for Point {
//     fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
//         write!(f, "({}, {})", self.x, self.y)
//     }
// }

// // get user full name
// impl User {
//     fn full_name(&self) -> String {
//         return format!("{} {}", self.first_name, self.last_name);
//     }
// }

// impl User {
//     fn age(&self) -> String {
//         return format!("{}", self.age);
//     }
// }

#[macro_use]
extern crate serde_derive;

use futures::executor::block_on;
use openapi::models::user;

extern crate reqwest;
extern crate serde;
extern crate serde_json;
extern crate url;
use std::time::{SystemTime, UNIX_EPOCH};

pub mod apis;
pub mod models;

#[tokio::main]
async fn main() {
    println!("Hello, world!");
    let conf = apis::configuration::Configuration::new();
    let existing_user = models::User::new(
        "test".to_string(),
        "test".to_string(),
        "test".to_string(),
        "email".to_string(),
    );

    let username = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_millis()
        .to_string();

    let email = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_millis()
        .to_string();

    let password = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_millis()
        .to_string();

    let req = models::CreateUserRequest::new(
        username,
        password,
        "first".to_string(),
        "last".to_string(),
        email,
    );

    let res = apis::user_api::create_user(&conf, req);
    // print res
    let user = res.await.unwrap();
    println!("response: {:?}", user);
}
