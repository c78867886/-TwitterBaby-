import { FollowerlistComponent } from './followerlist.component';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { DataService } from '../../services/data.service';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { DebugElement }    from '@angular/core';
import { By }              from '@angular/platform-browser';
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

describe('FollowerlistComponent', () => {
  let component: FollowerlistComponent;
  let fixture: ComponentFixture<FollowerlistComponent>;
  let spy: jasmine.Spy;
  let dataService: DataService;
  let de: DebugElement;
  let deArray: DebugElement[];
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
      declarations: [ FollowerlistComponent ],
      providers: [ {provide: 'data', useClass: DataService} ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FollowerlistComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get('data');
    spy = spyOn(dataService, 'showFollower')
          .and.returnValue(Promise.resolve([]));
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should have one follower in the list', () => {
    component.followerList = [{bio: 'test', username: 'test', id: 'test'}];
    fixture.detectChanges();
    de = fixture.debugElement.query(By.css('.flBio'));
    el = de.nativeElement;
    expect(el.textContent).toContain('test');
  });

  it('should have two follower in the list', () => {
    component.followerList = [{bio: 'test', username: 'test', id: 'test'},
                              {bio: 'test2', username: 'test2', id: 'test2'}];
    fixture.detectChanges();
    deArray = fixture.debugElement.queryAll(By.css('.flBio'));
    el = deArray[1].nativeElement;
    expect(el.textContent).toContain('test2');
  });
});
