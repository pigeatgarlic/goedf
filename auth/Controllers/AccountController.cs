using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Authenticator.Interfaces;
using Authenticator.Models.Auth;
using Authenticator.Models.User;
using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Authenticator.Database.Seeding;
using System.Linq;
using Authenticator.Database.Context;
using Authenticator.Models;
using Microsoft.Extensions.Options;
using Authenticator.Middleware;

namespace Authenticator.Controllers
{
    [ApiController]
    [Route("/account")]
    public class AccountController : ControllerBase
    {
        private readonly UserManager<UserAccount> _userManager;
        private readonly SignInManager<UserAccount> _signInManager;
        private readonly ITokenGenerator _tokenGenerator;
        private readonly GlobalDbContext _db;
        private readonly SystemConfig _config;

        public AccountController(
            UserManager<UserAccount> userManager,
            SignInManager<UserAccount> signInManager,
            ITokenGenerator tokenGenerator,
            GlobalDbContext db,
            IOptions<SystemConfig> config)
        {
            _config = config.Value;
            _userManager = userManager;
            _signInManager = signInManager;
            _tokenGenerator = tokenGenerator;
            _db = db;
        }

        [HttpPost]
        [Route("login")]
        public async Task<LoginResponse> Login([FromBody] LoginModel model)
        {
            if (!ModelState.IsValid)
            {
                return LoginResponse.GenerateFailure(model.UserName, new IdentityError {
                    Code = "Models Error",
                    Description = "Invalid Register Models" });
            }

            var result = await _signInManager.PasswordSignInAsync(model.UserName, model.Password, true, false);
            if (result.Succeeded)
            {
                UserAccount user = await _userManager.FindByNameAsync(model.UserName);
                string token = await _tokenGenerator.GenerateUserJwt(user);
                return LoginResponse.GenerateSuccessful(model.UserName, token, DateTime.Now.AddHours(1));
            }
            else
            {
                var error =  new List<IdentityError> ();
                error.Add(new IdentityError{
                    Code = "Login fail",
                    Description = "Wrong user or password" 
                });
                return LoginResponse.GenerateFailure(model.UserName, error);
            }
        }



        [HttpPost]
        [Route("register")]
        public async Task<LoginResponse> Register([FromBody] RegisterModel model)
        {

            if (!ModelState.IsValid)
            {
                return LoginResponse.GenerateFailure(model.Email, new IdentityError {
                    Code = "Models Error",
                    Description = "Invalid Register Models" });
            }

            var user = new UserAccount() {
                UserName = model.UserName,
                Email = model.Email,
                FullName = model.FullName,
                PhoneNumber = model.PhoneNumber,
            };
            if(model.DateOfBirth != null)
            {
                user.DateOfBirth = model.DateOfBirth;
            }
            if(model.Jobs != null)
            {
                user.Jobs = model.Jobs;
            }

            UserAccount userWithEmail = await _userManager.FindByEmailAsync(model.Email);
            if(userWithEmail != null)
            {
                var errors = new List<IdentityError>();
                errors.Add(new IdentityError{
                    Code = "Invalid email",
                    Description = "This email has been registered as an account"
                });

                return new LoginResponse {
                    Errors = errors,
                    UserName = model.UserName,
                };
            }

            var result = await _userManager.CreateAsync(user, model.Password);
            if (result.Succeeded)
            {
                UserAccount u = await _userManager.FindByEmailAsync(model.Email);
                await _userManager.AddToRoleAsync(u, RoleSeeding.USER);
                string token = await _tokenGenerator.GenerateUserJwt(u);
                return LoginResponse.GenerateSuccessful(model.UserName, token, DateTime.Now);
            }
            else
            {
                return LoginResponse.GenerateFailure(model.Email,result.Errors );
            }
        }









        [User]
        [HttpGet("infor")]
        public async Task<IActionResult> GetInfor()
        {
            var UserID = HttpContext.Items["UserID"];
            var account = await _userManager.FindByIdAsync(UserID.ToString());
            return Ok(new  UserInforModel
            {
                UserName = account.UserName,
                FullName = account.FullName,
                Jobs = account.Jobs,
                PhoneNumber = account.PhoneNumber,
                Gender = account.Gender,
                DateOfBirth = account.DateOfBirth,
                Avatar = account.Avatar
            });
        }

        [User]
        [HttpPost("infor")]
        public async Task<LoginResponse> SetAccountInfor([FromBody] UserInforModel infor)
        {
            var UserID = HttpContext.Items["UserID"];
            var account = await _userManager.FindByIdAsync(UserID.ToString());

            if(infor.Avatar != null)
            {
                account.Avatar = infor.Avatar;
            }
            if(infor.DateOfBirth != null)
            {
                account.DateOfBirth = infor.DateOfBirth;
            }
            if(infor.FullName != null)
            {
                account.FullName = infor.FullName;
            }
            if(infor.Gender != null)
            {
                account.Gender = infor.Gender;
            }
            if(infor.Jobs != null)
            {
                account.Jobs = infor.Jobs;
            }
            if(infor.PhoneNumber != null)
            {
                account.PhoneNumber = infor.Jobs;
            }
            if(infor.UserName != null)
            {
                account.UserName = infor.UserName;
            }
            
            var result = await _userManager.UpdateAsync(account);
            return result.Succeeded ? 
                LoginResponse.GenerateSuccessful(account.UserName,null,null) : 
                LoginResponse.GenerateFailure(account.UserName,result.Errors);
        }





        [User]
        [HttpPut("password")]
        public async Task<LoginResponse> UserGetRoles([FromBody] UpdatePasswordModel model)
        {
            IdentityResult result;
            var UserID = HttpContext.Items["UserID"];
            var account = await _userManager.FindByIdAsync(UserID.ToString());
            var hasPassword = await _userManager.HasPasswordAsync(account);

            if(hasPassword)
                result = await _userManager.ChangePasswordAsync(account,model.Old,model.New);
            else
                result = await _userManager.AddPasswordAsync(account,model.New);

            return result.Succeeded ? 
                LoginResponse.GenerateSuccessful(account.UserName,null,null) : 
                LoginResponse.GenerateFailure(account.UserName,result.Errors);
        }
    }
}
