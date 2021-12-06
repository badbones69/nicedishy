import React, { useEffect, useState } from 'react';
import { Utilities } from "../../utils/utilities";
import { useRouter } from 'next/router'
import Layout from "../../components/layout";

export default function Page() {
  const router = useRouter();

  return (
    <>
      <div className="container">
        <h2>Download the NiceDishy software</h2>
        <div className="row">
          <div className="col-6">
            <p>
              In order to participate in the NiceDishy network, you'll need to
              download and install the NiceDishy agent on a computer in your network.
              This computer should be able to connect to the dish, and be running 24x7.
            </p>
            <p>
              To get started, pick the right download on the right, and install.
              Once installed, you'll be prompted to log in and confirm which
              dish you are connecting.
            </p>
          </div>
          <dic className="col-1" />
          <div className="col-5">
            <div className="card" style={{width: "100%"}}>
              <div className="card-body">
                <h3><i class="bi bi-apple"></i> MacOS</h3>
                <a className="btn btn-outline-secondary" href="https://github.com/marc-campbell/nicedishy-macos/releases"><i class="bi bi-cloud-download"></i> Download</a>
                {' '}
                <a className="btn btn-outline-secondary" href="https://github.com/marc-campbell/nicedishy-macos"><i class="bi bi-github"></i> View Source Code</a>
                <br /><br /><strong>Current Version:</strong> 0.2
              </div>
            </div>
            <div className="card" style={{width: "100%"}}>
              <div className="card-body">
                <h3><i class="bi bi-windows"></i> Windows</h3>
                <p><em>Coming soon</em></p>
              </div>
            </div>
            <div className="card" style={{width: "100%"}}>
              <div className="card-body">
                <h3>Linux</h3>
                <p><em>Interested? Let us know.</em></p>
              </div>
            </div>
          </div>
        </div>
      </div>
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
