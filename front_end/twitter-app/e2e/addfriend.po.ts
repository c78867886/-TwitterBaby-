import { browser, by, element, protractor } from 'protractor';

export class FriendshipPage {
  sleep() {
    browser.sleep(500);
  }

  navigateToHome() {
    return browser.get('/home');
  }

  getFollowButton() {
    return element(by.css('.followbtn'));
  }

  getUserFollowing() {
    return element(by.css('.following'));
  }

  getFriendCol() {
    return element(by.css('mat-card'));
  }

  getNavRightButton() {
    return element(by.css('.navRightButton'));
  }

  getProfileBtn() {
    return element(by.css('.profileBtn'));
  }

  getProfileCpnt() {
    return element(by.css('.profileContainer'));
  }

  getLogoutButton() {
    return element(by.css('.logoutBtn'));
  }

}
