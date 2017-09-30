import { Injectable } from '@angular/core';
import { Tweet } from '../models/tweet.model';
import { Http, Response, Headers, RequestOptions } from '@angular/http';
import 'rxjs/add/operator/toPromise';
@Injectable()
export class DataService {
  localhost = "http://localhost:1323";
  constructor(private http: Http) { }

  getHeader(): RequestOptions {
    let access_token: string = localStorage.getItem("access_token");
    let headers: Headers = new Headers();
    headers.append('Authorization', 'Bearer ' + access_token);
    return new RequestOptions({ headers: headers });
  }

  getTweetList(id: string): Promise<Tweet[]> {
    let options: RequestOptions = this.getHeader();
    return this.http.get(this.localhost +`/api/v1/tweetlist/${id}`, options)
                      .toPromise()
                      .then((res: Response) => res.json())
                      .catch(this.handleError);
  }

  getUserInfo(id: string): Promise<Object> {
    let options: RequestOptions = this.getHeader();
    return this.http.get(this.localhost + `/api/v1/userInfo?username=${id}`, options)
                      .toPromise()
                      .then((res: Response) => res.json())
                      .catch(this.handleError);
  }

  followUser(mongoid: string): Promise<Object> {
    let post = {};
    let options: RequestOptions = this.getHeader();
    return this.http.post(this.localhost + `/api/v1/follow/${mongoid}`, post, options)
      .toPromise();
  }

  unfollowUser(mongoid: string) {
    
  }

  postTweet(text: string): Promise<object> {
    let post = {message: text};
    let options: RequestOptions = this.getHeader();
    return this.http.post(this.localhost + "/api/v1/newTweet", post, options)
      .toPromise();
  }

  mockLogin(): Promise<Object> {
    let loginfo: object = {email:"hojason117@gmail.com", password:"test1"};
    let headers: Headers = new Headers({ 'content-type': 'application/json'});
    let options: RequestOptions = new RequestOptions({ headers: headers });
    return this.http.post(this.localhost + '/api/v1/login', loginfo, options)
      .toPromise()
      .then((res: Response) => {
        console.log(res.json());
        localStorage.setItem('access_token', res.json().token);
        localStorage.setItem('id', res.json().id);
        localStorage.setItem('user_info_object', JSON.stringify(res.json()));
      })
      .catch(this.handleError);
  }

  // ERROR handler
  private handleError(error: any): Promise<any> {
    console.error('An error occurred', error);
    return Promise.reject(error.body || error);
  }
}
