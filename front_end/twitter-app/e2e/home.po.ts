import { browser, by, element, protractor } from 'protractor';

export class HomePage {
  // Async
  sleep() {
    browser.sleep(500);
  }

  // Navigate
  navigateToHome() {
    return browser.get('/home');
  }

  //Test for Home Page
  getNavBarBrand() {
    return element(by.css('.brand')).getText();
  }

  getPostArea() {
    return element(by.css('.expansionpanel'));
  }

  getPostInputArea() {
    return element(by.css('textarea'));
  }

  getPostSubmit() {
    return element(by.css('.submitbtn'));
  }

  getNewTweet() {
    return element.all(by.css('.tweetMessage')).first();
  }
  
  // Comment part
  getCommentBtn() {
    return element.all(by.css('#comment')).first();
  }

  getCommentText() {
    return element(by.css('.expansionpanel2'));
  }

  getCommentTextArea() {
    return element(by.css('.commenttextarea'));
  }

  getCommentSubmit() {
    return element(by.css('.actBtn2'));
  }

  getCommentList() {
    return element.all(by.css('.commentContent')).first();
  }

  getCancelBtn() {
    return element(by.css("#closeBtn"));
  }

  // Share part
  getShareBtn() {
    return element.all(by.css("#share")).first();
  }

  getShareInput() {
    return element(by.css(".shareInput"));
  }

  getShareSubmitbtn() {
    return element(by.css(".submitBtn"));
  }
}
