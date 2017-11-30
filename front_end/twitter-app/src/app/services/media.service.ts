import { Injectable } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Rx';
import 'rxjs/add/operator/toPromise';
import { User } from '../components/userprofile/user';

@Injectable()
export class MediaService {
  private backEndHostUrl: String = "http://127.0.0.1:1323/api/v1";

  constructor(private _http: Http) {  }

  // data = {
  //   size: '125422',
  //   type: 'image/jpeg',
  //   name:'test.jpg',
  //   url: base64,
  // };

  /**
   * Get header
   */
  getHeader(): RequestOptions {
    let access_token: string = localStorage.getItem("access_token");
    let headers: Headers = new Headers({'Content-Type': 'application/json'});
    headers.append('Authorization', 'Bearer ' + access_token);
    return new RequestOptions({ headers: headers });
  }


  uploadAvator(user: User) : Observable<object> {
    let options: RequestOptions = this.getHeader();
    let message: object = { picture: user.picture };
    console.log("picture data: " + user.picture);
    return this._http.post(`${this.backEndHostUrl}/updateProfilePic`, message ,options)
    .map(res => res.json())
    .do(res => {
       console.log("---------");
       console.log(res);
       console.log("Get result successfully!");
    })
    .catch(this.handleError);
  }

    // ERROR handler
    private handleError(error: any): Promise<any> {
      console.error('An error occurred', error);
      return Promise.reject(error.body || error);
    }

}
