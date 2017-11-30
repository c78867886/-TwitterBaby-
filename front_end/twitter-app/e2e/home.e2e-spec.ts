import { HomePage } from './home.po';

describe('twitter-app Home', () => {
  let page: HomePage;

  beforeEach(() => {
    page = new HomePage();
  });

  it('should post a new tweet', () => {
    page.getPostArea().click();
    expect(page.getPostSubmit().getAttribute('disabled')).toBe('true');
    page.sleep();

    page.getPostInputArea().sendKeys('Test Message for E2E test!');
    expect(page.getPostInputArea().getAttribute('value')).toEqual('Test Message for E2E test!');
    expect(page.getPostSubmit().getAttribute('disabled')).toBe(null);
    page.sleep();

    page.getPostSubmit().click();
    page.sleep();

    expect(page.getNewTweet().getText()).toEqual('Test Message for E2E test!');
  });

  it('should write a comment', () => {
    expect(page.getCommentBtn()).toBeTruthy();
    page.getCommentBtn().click();
    page.sleep();
    
    page.getCommentText().click();
    page.getCommentTextArea().sendKeys('Test Comment for E2E test!');
    expect(page.getCommentTextArea().getAttribute('value')).toEqual('Test Comment for E2E test!');
    page.getCommentSubmit().click();
    page.sleep();
    expect(page.getCommentList().getText()).toEqual('Test Comment for E2E test!');
    
    page.getCancelBtn().click();
    page.sleep();
  });

  it('should create a retweet', () => {
    expect(page.getShareBtn()).toBeTruthy();
    page.getShareBtn().click();
    page.sleep();
    page.getShareInput().sendKeys('This is a retweet message!');
    page.sleep();
    expect(page.getShareSubmitbtn()).toBeTruthy();
    page.getShareSubmitbtn().click();
    page.sleep();
  });

});
