
import { Component, OnInit, Input, Inject, ElementRef} from '@angular/core';
import { MatDialog } from '@angular/material'
import { EditUserProfileDialogComponent } from '../edit-user-profile-dialog/edit-user-profile-dialog.component';
import { Subscription } from 'rxjs/Subscription';
import {ChangeDetectorRef} from '@angular/core';
import { Router } from '@angular/router';
import { userinfo } from '../../models/userinfo.model';
import { MediaService } from '../../services/media.service';
import { User } from './user';
let URL: string = 'http://s7d2.scene7.com/is/image/PetSmart/PB1201_STORY_CARO-Authority-HealthyOutside-DOG-20160818?$PB1201$';
@Component({
  selector: 'app-userprofile',
  templateUrl: './userprofile.component.html',
  styleUrls: ['./userprofile.component.css']
})
export class UserprofileComponent implements OnInit {

  userName: string;
  userFirstName: string;
  userLastName: string;
  userEmail: string;
  userBio: string;
  dialogResult = "";
  avatar: any;
  private elem: ElementRef;
  userInfosubscription:Subscription;
  public file_srcs: string;
  public debug_size_before: string;
  public debug_size_after:string;
  user: User = new User();


  constructor(public dialog : MatDialog,
              private changeDetectorRef:ChangeDetectorRef,
              private route: Router,
              @Inject('media') private mediaService,
              @Inject('data') private data,
              ) { }

  ngOnInit() {
    let userinfo = JSON.parse(localStorage.getItem("user_info_object"));
    if (userinfo) {
        this.userName = userinfo.username;
        if (userinfo.picture !== "undefined"){
          console.log("user info has avator  " +userinfo.picture);
          this.avatar = userinfo.picture;
        } else {
          //this.avatar = 'assets/images/grumpycat.jpeg';
        }
    }

    /**
    * Get User Info
    */
    this.userInfosubscription = this.data.getUserInfoForProfile(this.userName)
    .subscribe(newUserInfo => 
      { 
        this.userFirstName = newUserInfo.userinfo.firstname;
        this.userLastName = newUserInfo.userinfo.lastname;
        this.userEmail = newUserInfo.userinfo.email;
        this.userBio = newUserInfo.userinfo.bio;
        this.avatar = newUserInfo.userinfo.picture;
        this.user.firstname = newUserInfo.userinfo.firstname;
        this.user.lastname = newUserInfo.userinfo.lastname;;
        this.user.email = newUserInfo.userinfo.email;
        this.user.bio = newUserInfo.userinfo.bio;
      }
    );
  }
  
  /**
   * Get User Info
   */
  getUserInfo() {
    let userinfo = JSON.parse(localStorage.getItem("user_info_object"));
    if (userinfo) {
        this.userName = userinfo.username;
        this.userFirstName = userinfo.firstname;
        this.userLastName = userinfo.lastname;
        this.userEmail = userinfo.email;
        this.userBio = userinfo.bio;
    }
  }

  /*
  * Edit profile
  * Open a dilog EditUserProfileDialogComponent
  */
  openEditUserProfileDiag() {
    let dialogRef = this.dialog.open(EditUserProfileDialogComponent, {
      width: '600px',
      data: 'The dialog data shows here',
    })

  /**
   *  After the dialog is closed
   */
  dialogRef.afterClosed().subscribe(result => {
    // Fresh user info 
    console.log('Dialog is closed: ${result}');
    this.dialogResult = result;
    let userinfo = JSON.parse(localStorage.getItem("user_info_object"));
    if (userinfo) {
        this.userName = userinfo.username;
        this.userFirstName = userinfo.firstname;
        this.userLastName = userinfo.lastname;
        this.userEmail = userinfo.email;
        this.userBio = userinfo.bio;
    }
    this.route.navigateByUrl('/userprofile');
    console.log(userinfo);
  })
  }

  /**
   * Change user's avator
   * @param input 
   */
  fileChange(input:any){
    if (input.target.files && input.target.files[0]){
      const reader = new FileReader();
      reader.readAsDataURL(input.target.files[0]);
      reader.onload=(x: any) =>{
        URL = x.target.result;
        this.avatar = URL;
        this.updateAvator();
      }
    }
  }

  /**
   * call media servers update Avator function
   */
  updateAvator(): void {
    console.log("start to request upload image");
    this.user.picture = this.avatar;
    this.mediaService.uploadAvator(this.user).subscribe(res =>{
    console.log("Upload successfully!");
    })
  }
}
