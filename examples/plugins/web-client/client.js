/**
 * FinFocus Connect Client - Browser Example
 *
 * This file demonstrates how to call FinFocus plugins using the Connect protocol.
 * The Connect protocol uses JSON over HTTP, making it trivial to call from browsers
 * without any special client libraries.
 *
 * Key benefits of Connect over pure gRPC:
 * - Works with standard fetch() API
 * - Uses JSON by default (human-readable)
 * - Works with HTTP/1.1 (no HTTP/2 requirement)
 * - Simple error handling with standard HTTP status codes
 */

// Get the server URL from the input field
function getServerUrl() {
    return document.getElementById('serverUrl').value.replace(/\/$/, '');
}

/**
 * Generic Connect RPC call helper
 *
 * @param {string} method - The RPC method name (e.g., 'Name', 'EstimateCost')
 * @param {object} request - The request payload (will be JSON-encoded)
 * @returns {Promise<object>} - The response payload
 */
async function connectRpc(method, request = {}) {
    const url = `${getServerUrl()}/finfocus.v1.CostSourceService/${method}`;

    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
    });

    // Check response.ok before parsing to handle non-JSON error responses
    // (e.g., HTML 502 from a reverse proxy)
    if (!response.ok) {
        let errorMessage = `Request failed with status ${response.status}`;
        try {
            // Try to parse Connect protocol error (JSON with 'code' and 'message' fields)
            const errorData = await response.json();
            if (errorData.message) {
                errorMessage = errorData.message;
            }
        } catch {
            // If JSON parsing fails, use the text response
            const text = await response.text();
            if (text) {
                errorMessage = text.substring(0, 200); // Limit error message length
            }
        }
        throw new Error(errorMessage);
    }

    return await response.json();
}

/**
 * Get the plugin name
 */
async function getName() {
    const resultDiv = document.getElementById('nameResult');
    resultDiv.style.display = 'block';
    resultDiv.className = 'result';
    resultDiv.textContent = 'Loading...';

    try {
        const response = await connectRpc('Name', {});
        resultDiv.className = 'result success';
        resultDiv.textContent = JSON.stringify(response, null, 2);
    } catch (error) {
        resultDiv.className = 'result error';
        resultDiv.textContent = `Error: ${error.message}`;
    }
}

/**
 * Estimate cost for a resource
 */
async function estimateCost() {
    const resultDiv = document.getElementById('estimateResult');
    resultDiv.style.display = 'block';
    resultDiv.className = 'result';
    resultDiv.textContent = 'Loading...';

    const resourceType = document.getElementById('resourceType').value;
    const instanceType = document.getElementById('instanceType').value;
    const region = document.getElementById('region').value;

    try {
        const response = await connectRpc('EstimateCost', {
            resource_type: resourceType,
            attributes: {
                instance_type: instanceType,
                region: region,
            },
        });
        resultDiv.className = 'result success';
        resultDiv.textContent = JSON.stringify(response, null, 2);

        // Format cost display
        if (response.cost_monthly !== undefined) {
            const formattedCost = new Intl.NumberFormat('en-US', {
                style: 'currency',
                currency: response.currency || 'USD',
            }).format(response.cost_monthly);
            resultDiv.textContent += `\n\nFormatted: ${formattedCost}/month`;
        }
    } catch (error) {
        resultDiv.className = 'result error';
        resultDiv.textContent = `Error: ${error.message}`;
    }
}

/**
 * Check if the plugin supports a resource type
 */
async function checkSupports() {
    const resultDiv = document.getElementById('supportsResult');
    resultDiv.style.display = 'block';
    resultDiv.className = 'result';
    resultDiv.textContent = 'Loading...';

    const provider = document.getElementById('supportsProvider').value;
    const resourceType = document.getElementById('supportsResourceType').value;

    // Note: The region field is required by the Supports RPC but is only used
    // for region-specific resource support checks. For general resource type
    // support queries, any valid region works. We use 'us-east-1' as a
    // reasonable default that exists across most cloud providers.
    try {
        const response = await connectRpc('Supports', {
            resource: {
                provider: provider,
                resource_type: resourceType,
                region: 'us-east-1',
            },
        });
        resultDiv.className = 'result success';
        resultDiv.textContent = JSON.stringify(response, null, 2);
    } catch (error) {
        resultDiv.className = 'result error';
        resultDiv.textContent = `Error: ${error.message}`;
    }
}

// Note: This file is designed for browser usage via <script> tag inclusion.
// Functions are exposed globally for use by index.html button onclick handlers.
