import { browser, by, element, protractor } from 'protractor';

export class LoginPage {
  // Async
  sleep() {
    browser.sleep(500);
  }
  
  // Navigate
  navigateToLogin() { 
    return browser.get('/login');
  }

  navigateToSign() {
    //browser.driver.manage().window().setSize(1280, 1024);
    return browser.get('/signup');
  }

  getParagraphText() {
    return element(by.css('h1')).getText();
  }


  // Test for Login Page
  loginUsername() {
    return element(by.css('input[name=username]'));
  }

  loginPassword() {
    return element(by.css('input[name=password]'));
  }

  loginSubmit() {
    browser.ignoreSynchronization = true;
    return element(by.css('form .btn'));
  }

  //Test for Signup page
  signUsername() {
    return element(by.css('input[name=username]'));
  }

  signEmail() {
    return element(by.css('input[name=email]'));
  }

  signFirstname() {
    return element(by.css('input[name=first_name]'));
  }

  signLastname() {
    return element(by.css('input[name=last_name]'));
  }

  signPassword() {
    return element(by.css('input[name=password]'));
  }

  signSubmitBtn() {
    return element(by.css('.btn'));
  }

  alertDiv() {
    return element(by.css('div .alert'));
  }
}
