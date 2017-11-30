import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Router } from '@angular/router';
import { FollowNotification } from '../../models/notification.model';
@Component({
  selector: 'app-nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.css']
})
export class NavBarComponent implements OnInit {
  username: string = "";
  loginName: string = "";
  hostName: string = "";
  shouldBeShowed: boolean;
  subscription: Subscription;
  rcvnewTweet: boolean = false;
  followNotificationList: Array<FollowNotification> = [];
  //connected: boolean = false;
  constructor(@Inject('data') private data, @Inject('auth') private auth, 
              @Inject('notify') private notify, private route: Router) { }

  ngOnInit() {
    
    this.notify.getEventListener().subscribe(event => {
      //console.log(this.notify.readyState());
      if (event.type === "open") {
        console.log("WS Connected!");
      } else if (event.type === "close") {
        console.log("WS Disconnected!");
      } else if (event.type === 'message') {
        console.log(event);
        if (event.data === "New tweets.") {
          this.rcvnewTweet = true;
        } else {
          this.followNotificationList.push(JSON.parse(event.data));
          console.log(this.followNotificationList);
        }
      }
    });

    this.subscription = this.auth.isLoggedIn().subscribe( bol => {
      if (bol) {
        this.shouldBeShowed = true;
        let userinfo = JSON.parse(localStorage.getItem("user_info_object"));
        this.loginName = userinfo.firstname + ' ' + userinfo.lastname;
        this.hostName = userinfo.username;
        if ( !this.notify.readyState() ) {
          this.notify.connect(userinfo.username);
        }
        
      } else {
        this.shouldBeShowed = false;
        
        if (this.notify.readyState()) {
          this.notify.close();
        }
        
      }
    });
  }
  
  ngOnDestroy() {
    if (this.notify.readyState()) {
      this.notify.close();
    }
    
  }

  onSubmit(): void {
    this.username = '';
    var input = document.getElementById("input");
    input.blur();
  }

  logout(): void {
    localStorage.clear();
    this.auth.isLoggedIn();
    this.notify.close();
    this.followNotificationList = [];
    this.route.navigateByUrl('/login');
  }

  refresh(): void {
    this.rcvnewTweet = false;
    if (this.route.url !== '/home') {
      this.route.navigateByUrl('/home')
        .catch(err => console.log(err));
    } else {
      this.data.getTweetListTimeLine(this.hostName, 1);
    }
    

  }

  sendClearMessage(): void {
    this.notify.send("Clear notifications.");
    this.followNotificationList = [];
  }

  /**
  * Go to user profile
  */
  goToUserProfile(): void{
    console.log("Navigate to the userprofile webpage");
    this.auth.isLoggedIn();
    this.route.navigateByUrl('/userprofile');
    console.log("Done");
  }
}
