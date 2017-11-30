import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { DataService } from '../../services/data.service';
import { PageSplitComponent } from './page-split.component';
import { HttpModule } from '@angular/http';
import { DebugElement }    from '@angular/core';
import { By }              from '@angular/platform-browser';
describe('PageSplitComponent', () => {
  let component: PageSplitComponent;
  let fixture: ComponentFixture<PageSplitComponent>;
  let de: DebugElement;
  let el: HTMLElement;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        HttpModule
      ],
      declarations: [ PageSplitComponent ],
      providers: [ {provide: 'data', useClass: DataService} ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PageSplitComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should create a 1 / 2 pagesplit', () => {
    component.totalPage = 2;
    component.index = 1;
    fixture.detectChanges();
    de = fixture.debugElement.query(By.css('p'));
    el = de.nativeElement;
    expect(el.textContent).toContain("1 / 2");
  });

  it('should add 1 for index', () => {
    component.totalPage = 2;
    component.index = 1;
    fixture.detectChanges();
    de = fixture.debugElement.query(By.css('p'));
    el = de.nativeElement;
    expect(el.textContent).toContain("1 / 2");
    component.index++;
    fixture.detectChanges();
    expect(el.textContent).toContain("2 / 2");
  });

  it('should reduce 1 for index', () => {
    component.totalPage = 2;
    component.index = 2;
    fixture.detectChanges();
    de = fixture.debugElement.query(By.css('p'));
    el = de.nativeElement;
    expect(el.textContent).toContain("2 / 2");
    component.index--;
    fixture.detectChanges();
    expect(el.textContent).toContain("1 / 2");
    component.nextPage();
    component.prePage();
  });
});
