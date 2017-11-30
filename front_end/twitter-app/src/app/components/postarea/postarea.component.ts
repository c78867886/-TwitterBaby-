import { Component, OnInit, Inject } from '@angular/core';
import { Router } from '@angular/router';
@Component({
  selector: 'app-postarea',
  templateUrl: './postarea.component.html',
  styleUrls: ['./postarea.component.css']
})
export class PostareaComponent implements OnInit {
  content: string = "";
  constructor(@Inject('data') private data,
              private router: Router) { }

  ngOnInit() {
  }

  postTweet(): void {
    this.data.postTweet(JSON.parse(localStorage.getItem("user_info_object")).username, this.content)
      .then(tweet => {
        this.content = "";
        console.log(tweet);
      })
      .catch(err => console.log(err));
  }

}
