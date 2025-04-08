import config from "@/config";

// Builds a URL by concatenating the API base URL with the given endpoint.
function buildUrl(endpoint) {
  return `${config.API_BASE_URL}${endpoint}`;
}

/**
 * Sends a POST request to the specified URL with the given body data.
 *
 * @param {string} url - The URL to send the request to.
 * @param {Object} bodyData - The data to send in the request body.
 * @returns {Promise<Object>} - A promise that resolves with the response data from the server.
 * @throws {Error|Object} - Throws an error object if the request fails.
 */
async function postRequest(url, bodyData) {
  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(bodyData),
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(JSON.stringify(data));
  }
  return data;
}

/**
 * Sends a request to the specified URL with the given method and body data.
 *
 * @param {string} url - The URL to send the request to.
 * @param {string} method - The HTTP method to use (e.g., "POST", "PUT", etc.).
 * @param {Object} bodyData - The data to send in the request body.
 * @returns {Promise<Object>} - A promise that resolves with the response data from the server.
 * @throws {Error|Object} - Throws an error object if the request fails.
 */
async function requestWithBody(url, method, bodyData) {
  try {
    const response = await fetch(url, {
      method: method.toUpperCase(), // PUT, POST, etc.
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(bodyData),
      credentials: "include",
    });
    const data = await response.json();
    if (!response.ok) {
      console.error(`API Error (${response.status}):`, data);
      throw { status: response.status, data: data };
    }
    return data;
  } catch (error) {
    console.error(`Request failed [${method} ${url}]:`, error);
    throw error;
  }
}

/**
 * Registers a new user with the given authentication data.
 *
 * @param {Object} authData - An object containing the authentication data
 *   (e.g., { email, password, username }).
 * @returns {Promise<Object>} - A promise that resolves with the response data
 *   from the server.
 * @throws {Error|Object} - Throws an error object if the request fails.
 */
export async function register(authData) {
  return await postRequest(buildUrl("/api/auth/register"), authData);
}

/**
 * Logs in an existing user with the given authentication data.
 *
 * @param {Object} authData - An object containing the authentication data
 *   (e.g., { email, password, username }).
 * @returns {Promise<Object>} - A promise that resolves with the response data
 *   from the server.
 * @throws {Error|Object} - Throws an error object if the request fails.
 */
export async function login(authData) {
  return await postRequest(buildUrl("/api/auth/login"), authData);
}

/**
 * Logs out the currently authenticated user and invalidates their
 * authentication cookie.
 *
 * @returns {Promise<void>}
 * @throws {Error} - Throws an error if the request fails.
 */
export async function logout() {
  return await postRequest(buildUrl("/api/auth/logout"), {});
}

/**
 * Fetches the profile of the currently authenticated user.
 *
 * @returns {Promise<Object>} - A promise that resolves to the user's profile data.
 * @throws {Error} - Throws an error if the request fails.
 */

export async function getUserProfile() {
  const response = await fetch(buildUrl("/api/users/me/profile"), {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(JSON.stringify(data));
  }
  return data;
}

/**
 * Sends a PUT request to update the currently authenticated user's profile.
 *
 * @param {Object} profileData - The profile data to update (e.g., { firstName, lastName }).
 * @returns {Promise<Object>} - A promise that resolves with the updated user data from the server.
 * @throws {Error|Object} - Throws an error object if the request fails.
 */
export async function updateUserProfile(profileData) {
  return await requestWithBody(
    buildUrl("/api/users/me/profile"),
    "PUT",
    profileData
  );
}

// TODO: Add a function to upload an avatar
// export async function uploadAvatar(formData) { ... }
