import { TestBed, inject } from '@angular/core/testing';
import { AuthService } from './auth.service';
import { HttpModule } from '@angular/http';
import { RouterTestingModule } from '@angular/router/testing';
import { AuthGuardService } from './auth-guard.service';

describe('AuthGuardService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpModule, RouterTestingModule],
      providers: [AuthGuardService,
      { provide: 'auth', useClass: AuthService }]
    });
  });
  
  it('should be created', inject([AuthGuardService], (service: AuthGuardService) => {
    expect(service).toBeTruthy();
    service.canActivate();
  }));
});
