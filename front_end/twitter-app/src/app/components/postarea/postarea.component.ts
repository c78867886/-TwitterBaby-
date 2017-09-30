import { Component, OnInit, Inject } from '@angular/core';

@Component({
  selector: 'app-postarea',
  templateUrl: './postarea.component.html',
  styleUrls: ['./postarea.component.css']
})
export class PostareaComponent implements OnInit {
  text: string = '';
  constructor(@Inject('data') private data) { }

  ngOnInit() {
  }

  postTweet(): void {
    if(this.text !== '') {
      this.data.postTweet(this.text);
      this.text = '';
    }  
  }

}
