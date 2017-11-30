import { Component, OnInit, Input, Inject } from '@angular/core';
import { Tweet } from '../../models/tweet.model';
import { EditCommentsDialogComponent } from '../edit-comments-dialog/edit-comments-dialog.component';
import { MatDialog } from '@angular/material'
import { Subscription } from 'rxjs/Subscription';
import { CommentlistComponent } from '../commentlist/commentlist.component';
import { RetweetDialogComponent } from '../retweet-dialog/retweet-dialog.component';

@Component({
  selector: 'app-tweetlist',
  templateUrl: './tweetlist.component.html',
  styleUrls: ['./tweetlist.component.css']
})
export class TweetlistComponent implements OnInit {
  @Input() tweetlist: Tweet[];
  @Input() username: string;
  yesShow: boolean;
  url: string = 'http://s7d2.scene7.com/is/image/PetSmart/PB1201_STORY_CARO-Authority-HealthyOutside-DOG-20160818?$PB1201$';
  constructor(public dialog : MatDialog, @Inject("data") private data) { }

  ngOnInit() {
  }

  editCommentDialog(tweet){
    console.log("Open a dialog");
    let dialogRef = this.dialog.open(EditCommentsDialogComponent, {
      width: '600px',
      data: tweet,
    })

    // After close the dialog
    dialogRef.afterClosed().subscribe(result => {
      console.log('Dialog is closed: ${result}');
    })

  }

  retweetDialog(tweet) {
    console.log("Click share");
    let dialogRef = this.dialog.open(RetweetDialogComponent, {
      width: '600px',
      data: tweet
    });
  }

  deleteTweet(tweetId) {
    this.data.deleteTweet(tweetId);
    this.tweetlist = this.tweetlist.filter((el) => {
      return el.id != tweetId;
    })
  }
}
