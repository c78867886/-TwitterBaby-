import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { DataService } from '../../services/data.service';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatToolbarModule, 
  MatInputModule, 
  MatMenuModule, 
  MatIconModule, 
  MatButtonModule, 
  MatCardModule, 
  MatExpansionModule,
  MatProgressSpinnerModule,
  MatChipsModule,
  } from '@angular/material';

import { UserInfoComponent } from './user-info.component';
import { UserInfo } from '../../models/userinfo.model';

describe('UserInfoComponent', () => {
  let component: UserInfoComponent;
  let fixture: ComponentFixture<UserInfoComponent>;
  let userInformation: UserInfo;
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        HttpModule,
        MatToolbarModule,
        MatInputModule,
        MatMenuModule,
        MatIconModule,
        MatButtonModule,
        MatCardModule,
        MatExpansionModule,
        MatProgressSpinnerModule,
        MatChipsModule,
        FormsModule,
        RouterTestingModule,
        BrowserAnimationsModule
      ],
      declarations: [ UserInfoComponent ],
      providers: [ {provide: 'data', useClass: DataService} ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserInfoComponent);
    component = fixture.componentInstance;
    userInformation = {followingcount: 1,
                      followercount: 2,
                      followed: false,
                      userinfo: {
                        bio: 'test',
                        email: 'test1',
                        firstname: "string;",
                        id: "string;",
                        lastname: "string;",
                        username: "string;"
                      }};
    component.userInfo = userInformation;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});