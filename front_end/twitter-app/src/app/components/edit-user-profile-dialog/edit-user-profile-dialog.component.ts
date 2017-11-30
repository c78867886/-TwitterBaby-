
import { Component, OnInit, Inject } from '@angular/core';
import { UserprofileComponent } from '../userprofile/userprofile.component';
import { MatDialogRef } from '@angular/material';
import { MAT_DIALOG_DATA } from '@angular/material';
import { DataService } from '../../services/data.service';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
import 'rxjs/add/operator/toPromise';


@Component({
  selector: 'app-edit-user-profile-dialog',
  templateUrl: './edit-user-profile-dialog.component.html',
  styleUrls: ['./edit-user-profile-dialog.component.css']
})
export class EditUserProfileDialogComponent implements OnInit {
  
  constructor(public thisDialogRef:MatDialogRef<EditUserProfileDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public dataDilog:string,
              @Inject('data') private data: DataService) {
               }

  userName: string;
  userFirstName: string;
  userLastName: string;
  userEmail: string;
  userBio: string;
  dialogResult = "";

  ngOnInit() {
    let userinfo = JSON.parse(localStorage.getItem("user_info_object"));
    if (userinfo) {
        this.userName = userinfo.username;
        this.userFirstName = userinfo.firstname;
        this.userLastName = userinfo.lastname;
        this.userEmail = userinfo.email;
        this.userBio = userinfo.bio;
    }
  }

  /**
   * Update User Profile
   */
  onCloseConfirm(){
      let userNewInfo: object = {"firstname": this.userFirstName,
                                "lastname": this.userLastName,
                                "bio":this.userBio,
                              };
      this.data.updateUserInfo(userNewInfo)
        .then(userinfo =>{
          console.log("Update successfully!");
        })
      
      this.thisDialogRef.close('Confirm');
  }

  /**
   * Cancel updating user profile
   */
  onCloseCancel(){
      this.thisDialogRef.close('Cancel');
  }
}
