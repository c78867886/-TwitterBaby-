import {Tweet} from "./tweet.model";
export interface UserInfo {
    followingcount: number;
    followercount: number;
    followed: boolean;
    userinfo: userinfo;
}


export interface userinfo {
    bio: string;
    email: string;
    firstname: string;
    id: string;
    lastname: string;
    username: string;
    picture ?: string;
}

export interface signUpUserInfo{
    username: string;
    firstname: string;
    lastname: string;
    email: string;
    passwd: string;
}