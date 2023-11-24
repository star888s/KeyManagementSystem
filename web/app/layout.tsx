import './globals.css';
import { Inter } from 'next/font/google';
import NextAuthProvider from './providers/NextAuth';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'KMS',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang='ja'>
      <body className={inter.className}>
        <NextAuthProvider>{children}</NextAuthProvider>
      </body>
    </html>
  );
}
