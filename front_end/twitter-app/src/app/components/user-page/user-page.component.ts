import { Component, OnInit, Inject } from '@angular/core';
import { Tweet } from '../../models/tweet.model';
import { ActivatedRoute, Router } from "@angular/router";
@Component({
  selector: 'app-user-page',
  templateUrl: './user-page.component.html',
  styleUrls: ['./user-page.component.css']
})
export class UserPageComponent implements OnInit {
  list: Tweet[];
  userInfo: object = null;
  username: string;
  isHost: boolean;
  page: number;
  totalPage: number;
  constructor(private route: ActivatedRoute,
    private _route: Router,
    @Inject('data') private data) { }

  ngOnInit() {
    this.route.params.subscribe(params => {
      
      this.data.getUserInfo(params["id"])
      .then(userinfo => 
        {
          console.log(userinfo);
          this.userInfo = userinfo;
          this.username = userinfo.userinfo.username;
          let userInfo = JSON.parse(localStorage.getItem("user_info_object"));
          this.isHost = this.username === userInfo.username ? true : false;
          this.username = userInfo.username;
        }     
      )
      .catch(err => {
        if(err.status === 404) {
          this._route.navigateByUrl('/404');
        }
      });;

      this.data.getTweetList(params["id"])
        .then(list => 
          {
            console.log(list);
            this.list = list.tweetlist;
            this.page = list.page;
            this.totalPage = list.totalPage;
          })
        .catch(err => {
          if(err.status === 404) {
            console.log("404");
            this.list = [];
          }
        });
      });
    }
}