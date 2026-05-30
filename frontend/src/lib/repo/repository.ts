import { PUBLIC_BASE_BE_URL } from "$env/static/public";
import axios from "axios";

const baseUrl = PUBLIC_BASE_BE_URL;

type clientCredentialsRespons = {
  clientKey: string;
  clientSecret: string;
  isSuccess: boolean;
};

export const fetchClientCredetials = async (
  token: string,
): Promise<clientCredentialsRespons> => {
  try {
    const res = await axios.get(`${baseUrl}/client/cs/${token}`);

    // logika ambil isi body
    const { clientKey, clientSecret } = res.data.data;

    if (!clientKey || !clientSecret) {
      console.log("clientKey or clientSecret is missing");
      return {
        clientKey: "",
        clientSecret: "",
        isSuccess: false,
      };
    }

    return {
      clientKey,
      clientSecret,
      isSuccess: true,
    };
  } catch (err) {
    if (axios.isAxiosError(err)) {
      console.error("Error message:", err.response?.data.message);
      return {
        clientKey: "",
        clientSecret: "",
        isSuccess: false,
      };
    } else {
      console.error("Unknown error: ", err);
      return {
        clientKey: "",
        clientSecret: "",
        isSuccess: false,
      };
    }
  }
};

type loginResponse = {
  isSuccess: boolean;
  jwtToken: string;
  message: string;
};

export const clientLogin = async (
  clientKey: string,
  clientSecret: string,
): Promise<loginResponse> => {
  try {
    const res = await axios.post(`${baseUrl}/auth/login`, {
      clientKey,
      clientSecret,
    });

    if (!res.data.jwtToken) {
      console.error("jwtToken is missing");
      return {
        isSuccess: false,
        jwtToken: "",
        message: "Token JWT di perlukan!",
      };
    }

    return {
      isSuccess: true,
      jwtToken: res.data.jwtToken,
      message: "Berhasil login!",
    };
  } catch (err) {
    if (axios.isAxiosError(err)) {
      console.error("Error message:", err.response?.data.message);
      return {
        isSuccess: false,
        jwtToken: "",
        message: err.response?.data.message ?? "Terjadi kesalahan!",
      };
    } else {
      console.error("Unknown error: ", err);
      return {
        isSuccess: false,
        jwtToken: "",
        message: "Terjadi kesalahan!",
      };
    }
  }
};
