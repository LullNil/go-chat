import config from "@/config";

// Builds a URL by concatenating the API base URL with the given endpoint.
function buildUrl(endpoint) {
  return `${config.API_BASE_URL}${endpoint}`;
}

/**
 * Makes an asynchronous POST request to the specified URL with the provided
 * authentication data.
 *
 * @param {string} url - The URL to which the request is sent.
 * @param {Object} authData - The authentication data to be sent in the
 * request body.
 *
 * @returns {Promise<Object>} - A promise that resolves with the JSON-parsed
 * response data if the request is successful.
 *
 * @throws {Error} - Throws an error if the server response is not successful,
 * including the error information from the server.
 */
async function request(url, authData) {
  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(authData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    // Throws error info from server
    throw new Error(JSON.stringify(data));
  }
  return data;
}

/**
 * Sends a login request to the server with the provided authentication data
 * and expects JWT-token.
 *
 * @param {Object} authData - The authentication data containing user
 * credentials.
 * @returns {Promise<Object>} - A promise that resolves with the server
 * response containing a token if successful.
 * @throws {Error} - Throws an error if the server response is not successful.
 */

export async function login(authData) {
  return await request(buildUrl("/login"), authData);
}

/**
 * Sends a registration request to the server with the provided authentication
 * data and expects a success response.
 *
 * @param {Object} authData - The authentication data containing
 * user credentials.
 * @returns {Promise<Object>} - A promise that resolves with the server
 * response if successful.
 * @throws {Error} - Throws an error if the server response is not successful.
 */
export async function register(authData) {
  return await request(buildUrl("/register"), authData);
}

/**
 * Sends a GET request to the server to obtain the data of the currently
 * authenticated user.
 *
 * @returns {Promise<Object>} - A promise that resolves with the server
 * response containing user data if successful.
 * @throws {Error} - Throws an error if the server response is not successful.
 */
export async function getUser() {
  const response = await fetch(buildUrl("/get-user"), {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    // Throws error info from server
    throw new Error(JSON.stringify(data));
  }
  return data;
}

/**
 * Sends a logout request to the server. This will invalidate the
 * authentication token on the server and clear the token stored on the
 * client.
 *
 * @returns {Promise<Object>} - A promise that resolves with the server
 * response if successful.
 * @throws {Error} - Throws an error if the server response is not successful.
 */
export async function logout() {
  const response = await fetch(buildUrl("/logout"), {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    // Throws error info from server
    throw new Error(JSON.stringify(data));
  }
  return data;
}
