import React, { useEffect, useState } from 'react';
import * as url from "url";
import { Utilities } from "../../utils/utilities";
import Layout from "../../components/layout";
import { useRouter } from 'next/router'
import Image from 'next/image';

export default function Page() {
  const router = useRouter();

  const onClickLogin = async (ev) => {
    ev.preventDefault();

    // ensure the user is logged out
    Utilities.logoutUser();

    // set or clear the next url from sessionstorage
    window.sessionStorage.removeItem('next');
    const query = url.parse(window.location.href, true).query;
    if (query && query.next) {
      window.sessionStorage.setItem('next', query.next);
    }

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_ENDPOINT}/api/v1/login?next=${encodeURIComponent(url.parse(window.location.href, true).query.next)}`);
      if (!res.ok) {
        console.log("error")
        return;
      }

      const data = await res.json();
      window.location.href = data.redirectURL;
    } catch(err) {
      console.log(err);
    }
  }

  useEffect( () => {
    if (Utilities.getToken()) {
      async function fetchData() {
        // validate that the token is still valid
        try {
          const res = await fetch(`${process.env.NEXT_PUBLIC_API_ENDPOINT}/api/v1/whoami`, {
          method: 'GET',
          headers: {
            "Content-Type": "application/json",
            "Authorization": Utilities.getToken(),
            },
          });

          if (!res.ok) {
            console.error("error");
            return;
          }

          const data = await res.json();
          if (data.user) {
            if (data.user.isWaitlisted) {
              router.push("/waitlist");
            } else {
              router.push("/dishies");
            }
            return;
          }
        } catch(err) {
          console.log(err);
        }
        Utilities.logoutUser();
      }
      fetchData();
    }
  })

  return (
    <>
      <h1>Log in to NiceDishy</h1>
      <strong>You will be taken to Google to authenticate.</strong>
      <p>
        By logging in, you are agreeing to our Terms of Service and Privacy Policy. We ask
        for read access to your Google profile in order to provide a complete experience
        here. We don&apos;t ask for permissions to change anything in your Google account.
      </p>
      <a href="#" width="80%" onClick={onClickLogin}>
        <Image src="/images/btn_google_signin_dark_pressed_web@2x.png" alt="Sign in with Google" width="200px" height="50px" />
      </a>
    </>
  );
}

Page.getLayout = function getLayout(page) {
  return (
    <Layout>
      {page}
    </Layout>
  );
}
