using Microsoft.AspNetCore.Identity;
using System;
using Microsoft.AspNetCore.Mvc;
using Authenticator.Interfaces;
using Authenticator.Models.User;
using Authenticator.Models.Auth;
using System.Threading.Tasks;
using Authenticator.Database.Seeding;
using Authenticator.Models;
using Google.Apis.Auth;
using Microsoft.Extensions.Options;

namespace Authenticator.Controllers
{
    [ApiController]
    [Route("/token")]
    public class TokenController : ControllerBase
    {
        private readonly UserManager<UserAccount> _userManager;
        private readonly SignInManager<UserAccount> _signInManager;
        private readonly ITokenGenerator _tokenGenerator;
        private readonly SystemConfig _config;

        public TokenController(
            UserManager<UserAccount> userManager,
            SignInManager<UserAccount> signInManager,
            ITokenGenerator tokenGenerator,
            IOptions<SystemConfig> config)
        {
            _userManager = userManager;
            _signInManager = signInManager;
            _tokenGenerator = tokenGenerator;
            _config = config.Value;
        }

        [HttpPost]
        [Route("challenge")]
        public async Task<ValidateResponse> Challene([FromBody] string token)
        {
            try {
                var account = await _tokenGenerator.ValidateUserToken(token);
                var roles = await _userManager.GetRolesAsync(account);
                var resp = new ValidateResponse { 
                    UserID = account.Id.ToString(),
                    Roles = await _tokenGenerator.GetUserRoles(account),
                    ValidatedAt = DateTime.Now
                };
                    
                return resp;
            }
            catch (Exception ex)
            {
                return new ValidateResponse {
                    UserID = null,
                    Roles = null,
                    ValidatedAt = DateTime.Now
                };
            }
        }

        [HttpPost]
        [Route("google")]
        public async Task<LoginResponse> Exchange([FromBody] string token)
        {
            var payload = await GoogleJsonWebSignature.ValidateAsync(token, new GoogleJsonWebSignature.ValidationSettings {
                Audience = new string[] { _config.GoogleOauthID }
            });

            var user = await _userManager.FindByEmailAsync(payload.Email);
            if (user == null)
            {
                user = new UserAccount
                {
                    UserName = payload.Email,
                    Email = payload.Email,
                    Avatar = payload.Picture,
                    FullName = payload.Name,
                    DateOfBirth = DateTime.Now,
                    PhoneNumber = "0123456789",
                    Jobs = "DefaultJob"
                };

                var result = await _userManager.CreateAsync(user);
                if (result.Succeeded)
                {
                    UserAccount u = await _userManager.FindByEmailAsync(payload.Email);
                    await _userManager.AddToRoleAsync(u, RoleSeeding.USER);
                }
            }

            // Add a login (i.e insert a row for the user in AspNetUserLogins table)
            await _signInManager.SignInAsync(user, isPersistent: false, "google");

            string ret = await _tokenGenerator.GenerateUserJwt(user);
            return  LoginResponse.GenerateSuccessful(user.UserName,ret,_tokenGenerator.GetTokenExpireTime());
        }
    }
}