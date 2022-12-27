export class User {
    created_at?: string;
    username: string;
    firstName: string;
    lastName: string;
    password: string;
    email: string;

    constructor(username: string, first_name: string, last_name: string, password: string, email: string) {
        this.firstName = first_name;
        this.lastName = last_name;
        this.username = username;
        this.password = password;
        this.email = email;
    }
}
