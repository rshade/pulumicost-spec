export class ValidationError extends Error {
  constructor(
    message: string,
    public readonly field?: string,
    public readonly code?: string
  ) {
    super(message);
    this.name = 'ValidationError';
    // Ensure instanceof works correctly
    Object.setPrototypeOf(this, ValidationError.prototype);
  }
}
