import { Injectable } from '@angular/core';
import { Tweet } from '../models/tweet.model';
import { Http, Response, Headers, RequestOptions } from '@angular/http';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Rx';
import 'rxjs/add/operator/toPromise';

@Injectable()
export class DataService {
  localhost = "http://localhost:1323";
  private timelineSource = new BehaviorSubject<object>([]);
  private backEndHostUrl: String = "http://127.0.0.1:1323/api/v1";
  constructor(private http: Http) { }

  // Create the header for http request
  getHeader(): RequestOptions {
    let access_token: string = localStorage.getItem("access_token");
    let headers: Headers = new Headers();
    headers.append('Authorization', 'Bearer ' + access_token);
    return new RequestOptions({ headers: headers });
  }


  // Get the tweetlist. This method will be changed to return Observable
  getTweetList(id: string): Promise<object> {
    let options: RequestOptions = this.getHeader();
    return this.http.get(this.localhost +`/api/v1/tweetlist/${id}?perpage=100&page=1`, options)
                      .toPromise()
                      .then((res: Response) => res.json())
                      .catch(this.handleError);
  }
  // Get the timeline of host.
  getTweetListTimeLine(id: string, page: number): Observable<object> {
    let options: RequestOptions = this.getHeader();
    this.timelineSource.next([]);
    this.http.get(this.localhost +`/api/v1/tweettimeline/${id}?perpage=15&page=${page}`, options)
                      .toPromise()
                      .then((res: Response) => {
                        console.log(res.json());
                        this.timelineSource.next(res.json());
                      })
                      .catch(this.handleError);
    return this.timelineSource.asObservable();
  }

  getUserInfo(id: string): Promise<Object> {
    let options: RequestOptions = this.getHeader();
    return this.http.get(this.localhost + `/api/v1/userInfo/${id}`, options)
                      .toPromise()
                      .then((res: Response) => res.json())
                      .catch(this.handleError);
  }
  
  // Follow and unfollow part
  followUser(id: string): Promise<Object> {
    let post = {};
    let options: RequestOptions = this.getHeader();
    return this.http.post(this.localhost + `/api/v1/follow/${id}`, post, options)
      .toPromise()
      .catch(this.handleError);
  }

  unfollowUser(id: string): Promise<Object> {
    let post = {};
    let options: RequestOptions = this.getHeader();
    return this.http.post(this.localhost + `/api/v1/unfollow/${id}`, post, options)
      .toPromise()
      .catch(this.handleError);
  }

  showFollower(id: string): Promise<object> {
    let options: RequestOptions = this.getHeader();
    return this.http.get(this.localhost + `/api/v1/showFollower/${id}`, options)
                      .toPromise()
                      .then((res: Response) => res.json())
                      .catch(this.handleError);
  }

  showFollowing(id: string): Promise<object> {
    let options: RequestOptions = this.getHeader();
    return this.http.get(this.localhost + `/api/v1/showFollowing/${id}`, options)
                      .toPromise()
                      .then((res: Response) => res.json())
                      .catch(this.handleError);
  }

  //Create new tweet
  postTweet(id: string, content: string): Promise<Object>{
    let options: RequestOptions = this.getHeader();
    let message: object = {message: content};
    //console.log(message);
    return this.http.post(this.localhost + `/api/v1/newTweet`, message, options)
            .toPromise()
            .then((res: Response) => {
              this.getTweetListTimeLine(id, 1);
              return res.json();
            })
            .catch(this.handleError);

  }

  //Delete Tweet
  deleteTweet(tweetId: string): Promise<Object> {
    let options: RequestOptions = this.getHeader();
    return this.http.delete(this.localhost + '/api/v1/deleteTweet/' + tweetId, options)
            .toPromise()
            .then((res: Response) => {
              return res.json();
            })
            .catch(this.handleError);
  }

  // ReTweet
  retweet(id: string, data: object): Promise<Object> {
    let options: RequestOptions = this.getHeader();
    return this.http.post(this.localhost + '/api/v1/reTweet', data, options)
            .toPromise()
            .then((res) => {
              this.getTweetListTimeLine(id, 1);
              return res.json();
            })
            .catch(this.handleError);
  }

  //MockLogin only for development
  // mockLogin(): Promise<Object> {
  //   let loginfo: object = {email:"hojason117@gmail.com", password:"test1"};
  //   let headers: Headers = new Headers({ 'content-type': 'application/json'});
  //   let options: RequestOptions = new RequestOptions({ headers: headers });
  //   return this.http.post('http://127.0.0.1:1323/api/v1/login', loginfo, options)
  //     .toPromise()
  //     .then((res: Response) => {
  //       console.log(res.json());
  //       localStorage.setItem('access_token', res.json().token);
  //       localStorage.setItem('id', res.json().id);
  //       localStorage.setItem('user_info_object', JSON.stringify(res.json()));
  //     })
  //     .catch(this.handleError);
  // }

  // ERROR handler
  private handleError(error: any): Promise<any> {
    console.error('An error occurred', error);
    return Promise.reject(error.body || error);
  }

  /**
   * Update User Info
   * By Diane
   * @param error 
   */
  updateUserInfo(userNewInfo): Promise<Object>{
    let options: RequestOptions = this.getHeader();
    // console.log(`${this.authUrl}/signup`, signUpUserInfo, options);
    return this.http.post(`${this.backEndHostUrl}/updateUserInfo`,userNewInfo, options)
      .toPromise()
      .then((res: Response) => {
        localStorage.setItem('user_info_object', JSON.stringify(res.json()));
        this.getUserInfoForProfile(userNewInfo.username);
        console.log(res.json());
    })
    .catch(this.handleError);
  }

  /**
   * Add a new comment into the tweet
   */
  addNewComment(commentContent, tweetid){
    let options: RequestOptions = this.getHeader();
    let message: object ={ "message":commentContent};
    return this.http.post(`${this.backEndHostUrl}/newcomment/${tweetid}`, message, options)
    .toPromise()
    .then((res: Response) => {
      console.log("back end response: successfully");
      console.log(JSON.stringify(res.json()));
      this.fetchComment(tweetid);
      return res.json();
      })
    .catch(this.handleError);
  }

  /**
   * Get comments for the specific tweet.
   * @param tweetid 
   */
  fetchComment(tweetid:string): Observable<object>{
    let options: RequestOptions = this.getHeader();
    return this.http.get(`${this.backEndHostUrl}/fetchcomment/${tweetid}`, options)
      .map(res => res.json())
      .do(res => {
         console.log(res);
         console.log("Get result successfully!");
      })
      .catch(this.handleError);
  }

  /**
   * Get User name by observable
   * @param username 
   */
  getUserInfoForProfile(username: string): Observable<object>{
    let options: RequestOptions = this.getHeader();
    return this.http.get(`${this.backEndHostUrl}/userInfo/${username}`, options)
      .map(res => res.json())
      .do(res => {
         console.log(res);
         console.log("Get user info successfully!");
      })
      .catch(this.handleError);
  }
}
