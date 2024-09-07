import { Link } from "react-router-dom";

function SignupCompleteMessage() {
  return (
    <section className="space-y-4 text-lg">
      <h2 className="mt-14 font-semibold text-xl">Thank you for siging up</h2>
      <p>
        We have sent a confirmation email to your provided email address. Click
        the provided link in the email to continue.
      </p>
      <p>
        <Link
          to="/login"
          className="font-semibold text-purple-500 hover:underline"
        >
          Login
        </Link>{" "}
        if you have already verified your email.
      </p>
    </section>
  );
}

export default SignupCompleteMessage;
