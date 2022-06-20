using Microsoft.AspNetCore.Identity;
using Microsoft.Extensions.Options;
using Microsoft.IdentityModel.Tokens;
using Authenticator.Interfaces;
using Authenticator.Models.Auth;
using Authenticator.Models.User;
using System;
using System.Collections.Generic;
using System.IdentityModel.Tokens.Jwt;
using System.Linq;
using System.Security.Claims;
using System.Text;
using System.Threading.Tasks;
using Authenticator.Database.Context;

namespace Authenticator.Services
{
    public class TokenGenerator : ITokenGenerator
    {

        private readonly JwtOptions _jwt;

        private readonly UserManager<UserAccount> _userManager;

        private readonly GlobalDbContext _db;

        public int ExpireHours {get;set;}

        public TokenGenerator(IOptions<JwtOptions> options, 
                              UserManager<UserAccount> userManager,
                              GlobalDbContext db )
        {
            _db = db;
            _jwt = options.Value;
            _userManager = userManager;
            ExpireHours = _jwt.ExpireHours;
        }

        public async Task<string> GenerateUserJwt(UserAccount user)
        {
            var userClaims = await _userManager.GetClaimsAsync(user);

            var tokenHandler = new JwtSecurityTokenHandler();
            var key = Encoding.ASCII.GetBytes(_jwt.Key);

            var tokenDescriptor = new SecurityTokenDescriptor
            {
                Subject = new ClaimsIdentity(new[] { new Claim("id", user.Id.ToString()) }),
                Expires = DateTime.Now.AddHours(_jwt.ExpireHours),
                SigningCredentials = new SigningCredentials(new SymmetricSecurityKey(key), SecurityAlgorithms.HmacSha256Signature),
            };

            var token = tokenHandler.CreateToken(tokenDescriptor);
            return tokenHandler.WriteToken(token);
        }

        public Task<UserAccount?> ValidateUserToken(string token)
        {
            try
            {
                var tokenHandler = new JwtSecurityTokenHandler();
                var key = Encoding.ASCII.GetBytes(_jwt.Key);
                tokenHandler.ValidateToken(token, new TokenValidationParameters
                {
                    ValidateIssuerSigningKey = true,
                    IssuerSigningKey = new SymmetricSecurityKey(key),
                    ValidateIssuer = false,
                    ValidateAudience = false,
                    RequireExpirationTime = true,
                }, out SecurityToken validatedToken);
                var jwtToken = (JwtSecurityToken)validatedToken;

                var id = jwtToken.Claims.First(x => x.Type == "id").Value;
                var account = _userManager.FindByIdAsync(id);
                return account;
            }
            catch 
            {
                return null;
            }
        }


        public async Task<List<string>> GetUserRoles(UserAccount account)
        {
            try
            {
                var ret = await _userManager.GetRolesAsync(account);
                return ret.ToList();
            }
            catch 
            {
                return null;
            }
        }

        public DateTime GetTokenExpireTime()  
        {
            return DateTime.Now.AddHours(ExpireHours);
        }
    }
}
