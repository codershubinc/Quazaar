const loginForm = document.getElementById('loginForm');
const errorMessage = document.getElementById('errorMessage');

loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    try {
        const response = await fetch('/api/v0.1/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // Save token
            localStorage.setItem('authToken', data.token);
            localStorage.setItem('username', data.username);
            
            // Redirect to main page
            window.location.href = '/';
        } else {
            showError(data.message || 'Login failed');
        }
    } catch (error) {
        showError('An error occurred. Please try again.');
        console.error('Login error:', error);
    }
});

function showError(message) {
    errorMessage.textContent = message;
    errorMessage.style.display = 'block';
}
