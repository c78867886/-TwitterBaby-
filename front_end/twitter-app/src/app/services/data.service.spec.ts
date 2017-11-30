import { TestBed, inject, getTestBed, fakeAsync, tick } from '@angular/core/testing';
import { MockBackend, MockConnection } from '@angular/http/testing';
import { Http, BaseRequestOptions, Response, ResponseOptions } from '@angular/http';
import { DataService } from './data.service';
import { HttpModule } from '@angular/http';
describe('DataService', () => {
  let backend: MockBackend;
  let service: DataService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [DataService, MockBackend, BaseRequestOptions,
      { provide: Http, useFactory: (backendInstance: MockBackend, defaultOptions: BaseRequestOptions) => {
                return new Http(backendInstance, defaultOptions);}, deps: [MockBackend, BaseRequestOptions] }],
      imports: [HttpModule]
    });
    backend = TestBed.get(MockBackend);
    service = TestBed.get(DataService);
  });

  it('should be created', inject([DataService], (service: DataService, mockBackend: MockBackend) => {
    expect(service).toBeTruthy();
  }));

  it('getTweetList() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.getTweetList("test")).toBeDefined();
  }));

  it('getTweetListTimeLine() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.getTweetListTimeLine("id", 1)).toBeDefined();
  }));

  it('getUserInfo() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.getUserInfo("id")).toBeDefined();
  }));

  it('getTweetList() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.getTweetList("test")).toBeDefined();
  }));

  it('followUser() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.followUser("id")).toBeDefined();
  }));

  it('unfollowUser() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.unfollowUser("id")).toBeDefined();
  }));

  it('showFollower() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.showFollower("id")).toBeDefined();
  }));

  it('showFollowing() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.showFollowing("id")).toBeDefined();
  }));

  it('postTweet() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.postTweet("id", "content")).toBeDefined();
  }));

  it('deleteTweet() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.deleteTweet("id")).toBeDefined();
  }));

  it('retweet() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.retweet("id", {})).toBeDefined();
  }));

  it('updateuserinfo() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.updateUserInfo({})).toBeDefined();
  }));

  it('addnewcomment() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.addNewComment({}, "222")).toBeDefined();
  }));

  it('fetchComment() should send request to server', fakeAsync(() => {
    backend.connections.subscribe(connection => { 
      let options = new ResponseOptions({
        body: JSON.stringify({success: true})
      });
      connection.mockRespond(new Response(options));
    });
    tick();
    expect(service.fetchComment("222")).toBeDefined();
  }));

});
