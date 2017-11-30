export class User {
    constructor(
        public _id: String = '',
        public username: string = '',
        public firstname: string = '',
        public lastname: string = '',
        public email: string = '',
        public picture: string = '',
        public bio: string = '',
    ) {}
}