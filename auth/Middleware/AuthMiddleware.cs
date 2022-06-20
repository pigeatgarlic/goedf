using Microsoft.AspNetCore.Http;
using System.Linq;
using Microsoft.Extensions.Options;
using System.Collections.Generic;
using System.Threading.Tasks;
using Newtonsoft.Json;
using Microsoft.AspNetCore.Http.Features;
using System.Net;
using RestSharp;
using System;
using Authenticator.Models;
using Authenticator.Models.User;
using Authenticator.Interfaces;
using Microsoft.AspNetCore.Identity;
using Authenticator.Database.Seeding;


namespace Authenticator.Middleware
{
    public class JwtMiddleware
    {
        private readonly RequestDelegate _next;


        private readonly string IssuerUrl;

        private readonly ITokenGenerator _token;

        public JwtMiddleware(RequestDelegate next, 
                            ITokenGenerator token,
                            IOptions<SystemConfig> config)
        {
            _next = next;
            _token = token;
        }

        public async Task Invoke(HttpContext context)
        {
            string token = context.Request.Headers["Authorization"].FirstOrDefault()?.Split(" ").Last();
            if (token != null)
            {
                await attachUserToContext(context,  token);
            }
            await _next(context);
        }

        private async Task attachUserToContext(HttpContext context, string token)
        {
            try
            {
                var result = await _token.ValidateUserToken(token);
                context.Items.Add("UserID", result.Id.ToString());
            }
            catch (Exception ex)
            {
                // do nothing if jwt validation fails
                // user is not attached to context so request won't have access to secure routes
            }
        }
    }



    public class AuthorizeMiddleWare
    {
        private readonly RequestDelegate _next;

        private readonly SystemConfig _config;

        private readonly ITokenGenerator _token;

        private readonly UserManager<UserAccount> _manager;

        public AuthorizeMiddleWare(RequestDelegate next, 
                                   ITokenGenerator token,
                                   UserManager<UserAccount> manager,
                                   IOptions<SystemConfig> config)
        {
            _next = next;
            _config = config.Value;
            _token = token;
            _manager = manager;
        }

        public async Task Invoke(HttpContext context)
        {
            bool next = false;
            List<string> roles = new List<string>();

            var endpoint = context.Features.Get<IEndpointFeature>()?.Endpoint;

            var userAttribute = endpoint?.Metadata.GetMetadata<UserAttribute>();
            var adminAttribute = endpoint?.Metadata.GetMetadata<AdminAttribute>();
            
            if (userAttribute == null && adminAttribute == null)
            {
                next = true;
            }
            else
            {
                string ID = (string)context.Items["UserID"];
                roles = await _token.GetUserRoles(await _manager.FindByIdAsync(ID));
            }

            if (userAttribute != null)
            {
                roles.ForEach(i => {
                    if (i == RoleSeeding.USER)
                    {
                        next = true;
                    }
                });
            }
            if (adminAttribute != null)
            {
                roles.ForEach(i => {
                    if (i == RoleSeeding.ADMIN)
                    {
                        next = true;
                    }
                });
            }


            if (next)
            {
                await _next(context);
            }
        }
    }
}