import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material';
import { MatDialogRef } from '@angular/material';
@Component({
  selector: 'app-retweet-dialog',
  templateUrl: './retweet-dialog.component.html',
  styleUrls: ['./retweet-dialog.component.css']
})
export class RetweetDialogComponent implements OnInit {
  message: string = "";
  constructor(@Inject(MAT_DIALOG_DATA) private tweet,
              private thisDialogRef: MatDialogRef<RetweetDialogComponent>,
              @Inject("data") private data) { }
  
  ngOnInit() {
  }

  onSubmit(tweetid: string) {
    let dataObject = {
      Idretweet: tweetid,
      message: this.message
    }
    this.data.retweet(JSON.parse(localStorage.getItem("user_info_object")).username, dataObject);
    this.thisDialogRef.close("Share!");
  }

}
