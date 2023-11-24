'use client';
import { signIn, useSession } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import React from 'react';

export default function Login() {
  const { data: session } = useSession();
  const router = useRouter();

  if (!session) {
    return <button onClick={() => signIn('cognito')}>サインイン</button>;
  }

  router.push('/home');

  return null;
}
