import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormsModule } from '@angular/forms';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { AppComponent } from './app.component';
import { DataService } from './services/data.service';
import { AuthService } from './services/auth.service';
import { AuthGuardService } from './services/auth-guard.service';
import { AuthGuardLoggedService } from './services/auth-guard-logged.service';
import { NotificationService } from './services/notification.service';
import { rooting } from './app.router';
import { MediaService } from './services/media.service';


import { MatToolbarModule, 
         MatInputModule, 
         MatMenuModule, 
         MatIconModule, 
         MatButtonModule, 
         MatCardModule, 
         MatExpansionModule,
         MatProgressSpinnerModule,
         MatChipsModule,
         MatTooltipModule,
         MatDialogModule
         } from '@angular/material';
         
import { UserInfoComponent } from './components/user-info/user-info.component';
import { TweetlistComponent } from './components/tweetlist/tweetlist.component';
import { PostareaComponent } from './components/postarea/postarea.component';
import { LoginComponent } from './components/login/login.component';
import { UserPageComponent } from './components/user-page/user-page.component';
import { FollowlistComponent } from './components/followlist/followlist.component';
import { FollowerlistComponent } from './components/followerlist/followerlist.component';
import { PageSplitComponent } from './components/page-split/page-split.component';
import { UserloginComponent } from './components/userlogin/userlogin.component';
import { SignUpComponent } from './components/sign-up/sign-up.component';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';
import { NotfoundComponent } from './components/notfound/notfound.component';
import { UserprofileComponent } from './components/userprofile/userprofile.component';
import { EditUserProfileDialogComponent } from './components/edit-user-profile-dialog/edit-user-profile-dialog.component';
import { EditCommentsDialogComponent } from './components/edit-comments-dialog/edit-comments-dialog.component';
import { CommentlistComponent } from './components/commentlist/commentlist.component';
import { RetweetDialogComponent } from './components/retweet-dialog/retweet-dialog.component';
import { FooterComponent } from './components/footer/footer.component';


@NgModule({
  declarations: [
    AppComponent,
    NavBarComponent,
    UserInfoComponent,
    TweetlistComponent,
    PostareaComponent,
    LoginComponent,
    UserPageComponent,
    FollowlistComponent,
    FollowerlistComponent,
    PageSplitComponent,
    UserloginComponent,
    SignUpComponent,
    NotfoundComponent,
    UserprofileComponent,
    EditUserProfileDialogComponent,
    EditCommentsDialogComponent,
    CommentlistComponent,
    RetweetDialogComponent,
    FooterComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpModule,
    FormsModule,
    MatToolbarModule,
    MatInputModule,
    MatMenuModule,
    MatIconModule,
    MatButtonModule,
    MatCardModule,
    MatExpansionModule,
    MatProgressSpinnerModule,
    MatChipsModule,
    MatTooltipModule,
    ReactiveFormsModule,
    MatDialogModule,
    rooting,
  ],
  entryComponents:[
    EditUserProfileDialogComponent,
    EditCommentsDialogComponent,
    RetweetDialogComponent,
  ],
  providers: [
    { provide: 'data', useClass: DataService },
    { provide: 'auth', useClass: AuthService },
    { provide: 'media', useClass: MediaService},
    { provide: "notify", useClass: NotificationService},
    AuthGuardService, AuthGuardLoggedService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
