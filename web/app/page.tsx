'use client';
import { useSession } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import React from 'react';
import { Login } from './components/AuthButton';

export default function Home() {
  const { data: session } = useSession();
  const router = useRouter();
  if (!session) {
    return (
      <main className='flex h-screen justify-center items-center'>
        <div className='text-center'>
          <h1 className='text-4xl font-bold'>KMS</h1>
          <Login></Login>
        </div>
      </main>
    );
  }
  router.push('/home');

  return null;
}
