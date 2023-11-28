import NextAuth from 'next-auth';
import Providers from 'next-auth/providers/cognito';

const handler = NextAuth({
  providers: [
    Providers({
      clientId: process.env.COGNITO_CLIENT_ID,
      clientSecret: process.env.COGNITO_CLIENT_SECRET,
      issuer: process.env.COGNITO_DOMAIN,
      checks: 'nonce',
    }),
  ],
  secret: process.env.SECRET,
});

export { handler as GET, handler as POST };
