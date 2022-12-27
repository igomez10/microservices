import type { AppProps } from 'next/app'
import { User } from '../models/user';
export default function UserConfirmation(user: User) {
    return (
        <div>
            <h1>Confirmation</h1>
            <p>Thank you for signing up!</p>
            <p>Username: {user.username}</p>
            <p>First Name: {user.firstName}</p>
            <p>Last Name: {user.lastName}</p>
            <p>Email: {user.email}</p>
            <p>Password: {user.password}</p>
            <p>Created At: {user.created_at}</p>
        </div>
    )
}
