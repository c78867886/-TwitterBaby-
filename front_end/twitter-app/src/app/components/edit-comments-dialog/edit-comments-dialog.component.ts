import { Component, OnInit, Input, Inject} from '@angular/core';
import { MatDialogRef } from '@angular/material';
import { MAT_DIALOG_DATA } from '@angular/material';
import { Subscription } from 'rxjs/Subscription';
import { Tweet } from '../../models/tweet.model';
import { CommentlistComponent } from '../commentlist/commentlist.component';

@Component({
  selector: 'app-edit-comments-dialog',
  templateUrl: './edit-comments-dialog.component.html',
  styleUrls: ['./edit-comments-dialog.component.css']
})
export class EditCommentsDialogComponent implements OnInit {

  public commentContent: string = '';
  commentsSubscription: Subscription;
  list:Comment[];
  url: string = 'http://s7d2.scene7.com/is/image/PetSmart/PB1201_STORY_CARO-Authority-HealthyOutside-DOG-20160818?$PB1201$';
  constructor(public thisDialogRef:MatDialogRef<EditCommentsDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public tweet,
              @Inject('data') private data,
             ) { }

  ngOnInit() {
    this.freshData();
  }

  /**
   * Add a comment
   */
  postComment(){
    this.data.addNewComment(this.commentContent, this.tweet.id)
    .then(res => {
      this.freshData();
      this.commentContent = '';
      console.log("Adding comment is done")
      console.log(res);
    })
  }

  /**
   * Fresh Comment list
   */
  private freshData(): void{
    this.commentsSubscription = this.data.fetchComment(this.tweet.id)
      .subscribe(list => 
        {
          this.list = list.commentlist;
    })
  }

  onCloseSubmit(){
    this.data.addNewComment(this.commentContent, this.tweet.id)
      .then(res => {
        console.log("Adding comment is done")
        console.log(res);
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
