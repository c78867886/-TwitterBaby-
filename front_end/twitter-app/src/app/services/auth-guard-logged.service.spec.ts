import { TestBed, inject } from '@angular/core/testing';
import { AuthService } from './auth.service';
import { HttpModule } from '@angular/http';
import { RouterTestingModule } from '@angular/router/testing';
import { AuthGuardLoggedService } from './auth-guard-logged.service';

describe('AuthGuardLoggedService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpModule, RouterTestingModule],
      providers: [AuthGuardLoggedService,
      { provide: 'auth', useClass: AuthService }]
    });
  });

  it('should be created', inject([AuthGuardLoggedService], (service: AuthGuardLoggedService) => {
    expect(service).toBeTruthy();
    service.canActivate();
  }));
});
