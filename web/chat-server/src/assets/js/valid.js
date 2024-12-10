export function checkTelephoneValid(telephone) {
    const regex = /^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$/;
    return regex.test(telephone);
};

export function checkEmailValid(email) {
    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return regex.test(email);
};

