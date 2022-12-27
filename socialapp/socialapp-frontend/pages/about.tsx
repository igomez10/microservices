import React, { useState } from 'react'
import { User } from '../models/user';
import UserConfirmation from './user_confirmation';


export default function About() {
    const [firstName, setFirstName] = useState("test" + Math.random());
    const [lastName, setLastName] = useState("test" + Math.random());
    const [username, setUsername] = useState("test" + Math.random());
    const [password, setPassword] = useState("test" + Math.random());
    const [email, setEmail] = useState("test" + Math.random());
    const [createdUser, setCreatedUser] = useState(new User("", "", "", "", ""));
    // let createdUser: User = new User("first", "last", "", "", "");



    // create user function returns a user object
    let createUser = async (user: User): Promise<User> => {
        console.log("Creating user");

        let res = await fetch('http://localhost:8085/users', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: user.username,
                first_name: user.firstName,
                last_name: user.lastName,
                password: user.password,
                email: user.email,
            }),
        })
            .then((response) => {
                let data = response.json()
                console.log('First:', data);
                return data;
            })
            .catch((error) => {
                console.error('Error:', error);
                return error;
            });
        return res;
    };



    return (
        <div>

            {/* Form for creating a user */}
            <form>
                <label htmlFor="first_name">First Name</label>
                <input type="text" id="first_name" name="first_name" value={firstName} placeholder="First Name" onChange={(e) => setFirstName(e.target.value)} />
                <label htmlFor="last_name">Last Name</label>
                <input type="text" id="last_name" name="last_name" value={lastName} placeholder="Last Name" onChange={(e) => setLastName(e.target.value)} />
                <label htmlFor="username">Username</label>
                <input type="text" id="username" name="username" value={username} placeholder="Username" onChange={(e) => setUsername(e.target.value)} />
                <label htmlFor="email">Email</label>
                <input type="text" id="email" name="email" value={email} placeholder="Email" onChange={(e) => setEmail(e.target.value)} />
                <label htmlFor="password">Password</label>
                <input type="text" id="password" name="password" value={password} placeholder="Password" onChange={(e) => setPassword(e.target.value)} />
                <button type="submit" onClick={async (e) => {
                    e.preventDefault();
                    let user = new User(username, firstName, lastName, password, email);
                    createUser(user).then((createdUser) => {
                        console.log("Created User: ", createdUser);
                        user.created_at = createdUser.created_at;
                        setCreatedUser(user);
                    }).catch((error) => {
                        console.log("Error: ", error);
                    });


                }}
                > Create User</button>
            </form>
            <UserConfirmation {...createdUser} />

        </div >
    )
}
