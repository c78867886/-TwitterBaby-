import { NavBarComponent } from './nav-bar.component';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { DebugElement }    from '@angular/core';
import { By }              from '@angular/platform-browser';
import { DataService } from '../../services/data.service';
import { NotificationService } from '../../services/notification.service';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { AuthService } from '../../services/auth.service';
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

describe('NavBarComponent', () => {
  let component: NavBarComponent;
  let fixture: ComponentFixture<NavBarComponent>;
  let de: DebugElement;
  let el: HTMLElement;

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
      declarations: [ NavBarComponent ],

      providers: [ {provide: 'data', useClass: DataService},
                   {provide: 'auth', useClass: AuthService},
                   {provide: 'notify', useClass: NotificationService} ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NavBarComponent);
    component = fixture.componentInstance;
    component.shouldBeShowed = true;
    fixture.detectChanges();
  });

  it('should be created', () => {
    component.shouldBeShowed = true;
    component.ngOnInit();
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should have a refresh button', () => {
    component.shouldBeShowed = true;
    fixture.detectChanges();
    de = fixture.debugElement.query(By.css('.refresh'));
    el = de.nativeElement;
    // el.click();
    component.refresh();
    expect(el).toBeTruthy();
  });
  
  it('should have a notification button', () => {
    component.shouldBeShowed = true;
    fixture.detectChanges();
    de = fixture.debugElement.query(By.css('.notifications'));
    el = de.nativeElement;
    //component.sendClearMessage();
    expect(el).toBeTruthy();
  });

  it('should have a user management button', () => {
    component.shouldBeShowed = true;
    fixture.detectChanges();
    de = fixture.debugElement.query(By.css('.navRightButton'));
    el = de.nativeElement;
    expect(el).toBeTruthy();
  });

  it('should open user management dialog', () => {
    component.shouldBeShowed = true;
    fixture.detectChanges();
    de = fixture.debugElement.query(By.css('.navRightButton'));
    el = de.nativeElement;
    el.click();
    expect(el).toBeTruthy();
    component.ngOnInit();
    component.ngOnDestroy();
    //component.sendClearMessage();
    component.goToUserProfile();
    //component.logout();
    component.onSubmit();
  });
});
