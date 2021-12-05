
export default function Footer() {
  return (
    <div className="container">
      <footer className="py-5">
        <div className="row">
          <div className="col-2">
            <h5>Dishy Owners</h5>
            <ul className="nav flex-column">
              <li className="nav-item mb-2"><a href="#" className="nav-link p-0 text-muted">Join NiceDishy</a></li>
              <li className="nav-item mb-2"><a href="#" className="nav-link p-0 text-muted">Compare My Dishy</a></li>
              <li className="nav-item mb-2"><a href="#" className="nav-link p-0 text-muted">Is Everything Ok?</a></li>
            </ul>
          </div>

          <div className="col-2">
            <h5>Dishy Waitlisters</h5>
            <ul className="nav flex-column">
              <li className="nav-item mb-2"><a href="#" className="nav-link p-0 text-muted">Get Prepared</a></li>
            </ul>
          </div>

          <div className="col-2">
            <h5>Everyone Else</h5>
            <ul className="nav flex-column">
              <li className="nav-item mb-2"><a href="#" className="nav-link p-0 text-muted">See Some Stats</a></li>
              <li className="nav-item mb-2"><a href="#" className="nav-link p-0 text-muted">Request Data</a></li>
            </ul>
          </div>

          <div className="col-4 offset-1">
            <form>
              <h5>Subscribe to our newsletter</h5>
              <p>A semi-regular newsletter about Dishy.</p>
              <div className="d-flex w-100 gap-2">
                <label for="newsletter1" className="visually-hidden">Email address</label>
                <input id="newsletter1" type="text" className="form-control" placeholder="Email address" />
                <button className="btn btn-primary" type="button">Subscribe</button>
              </div>
            </form>
          </div>
        </div>

        <div className="d-flex justify-content-between py-4 my-4 border-top">
          <p>&copy; 2021. All rights reserved.</p>
          <ul className="list-unstyled d-flex">

          </ul>
        </div>
      </footer>
    </div>
  );
}
