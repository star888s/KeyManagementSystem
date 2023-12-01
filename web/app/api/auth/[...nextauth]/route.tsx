import NextAuth from 'next-auth';
import Providers from 'next-auth/providers/cognito';

const handler = NextAuth({
  providers: [
    Providers({
      clientId: process.env.COGNITO_CLIENT_ID,
      clientSecret: process.env.COGNITO_CLIENT_SECRET,
      issuer: process.env.COGNITO_DOMAIN,
    }),
  ],
});

export { handler as GET, handler as POST };
